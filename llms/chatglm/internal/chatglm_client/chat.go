package chatglm_client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

const (
	// chatglm_pro  chatglm_std  chatglm_lite  characterglm(超拟人大模型)
	defaultChatModel = "chatglm_std"
	eventAdd         = "add"
	eventFinish      = "finish"
	eventError       = "error"
	eventnterrupted  = "interrupted"
)

type ChatRequest struct {
	Model       string         `json:"model,omitempty"`
	Prompt      []*ChatMessage `json:"prompt,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	TopP        float64        `json:"top_p,omitempty"`
	RequestId   string         `json:"request_id,omitempty"`
	// SSE接口调用时，用于控制每次返回内容方式是增量还是全量，不提供此参数时默认为增量返回, true 为增量返回 , false 为全量返回
	Incremental bool `json:"incremental,omitempty"`
	// sse返回需设置streamingFunc
	// 结束时返回一个错误 Return an error to stop streaming early.
	StreamingFunc func(ctx context.Context, chunk []byte) error `json:"-"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StreamedChatResponsePayload struct {
	ID    string `json:"ID,omitempty"`
	Event string `json:"event,omitempty"`
	Data  string `json:"data,omitempty"`
	Meta  Meta   `json:"meta,omitempty"`
}

type Meta struct {
	Usage Usage `json:"usage"`
}

type ChatResponse struct {
	Code    int    `json:"code,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

func (c *Client) createChat(ctx context.Context, payload *ChatRequest) (*ChatResponse, error) {
	var method string
	if payload.StreamingFunc != nil {
		method = "sse-invoke"
	} else {
		method = "invoke"
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.buildURL(c.model, method), body)
	if err != nil {
		return nil, err
	}
	c.setHeader(req)
	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("API returned unexpected status code: %d", r.StatusCode)

		// No need to check the error here: if it fails, we'll just return the
		// status code.
		var errResp errorMessage
		if err := json.NewDecoder(r.Body).Decode(&errResp); err != nil {
			return nil, errors.New(msg) // nolint:goerr113
		}

		return nil, fmt.Errorf("%s: %s", msg, errResp.Error.Message) // nolint:goerr113
	}
	if payload.StreamingFunc != nil {
		return parseStreamingChatResponse(ctx, r, payload)
	}
	var response ChatResponse
	return &response, json.NewDecoder(r.Body).Decode(&response)
}

func parseStreamingChatResponse(ctx context.Context, r *http.Response, payload *ChatRequest) (*ChatResponse, error) { //nolint:cyclop,lll
	scanner := bufio.NewScanner(r.Body)
	responseChan := make(chan StreamedChatResponsePayload)
	go func() {
		defer close(responseChan)
		for scanner.Scan() {
			line := scanner.Text()
			log.Println(line)
			if line == "" {
				continue
			}
			chunkResp, err := decodeStreamData(line)
			if err != nil {
				log.Printf("failed to decode stream payload: %v", err)
				break
			}
			responseChan <- *chunkResp
		}
		if err := scanner.Err(); err != nil {
			log.Println("issue scanning response:", err)
		}
	}()
	// Parse response
	response := ChatResponse{
		Data: Data{
			Choices: []*Choices{
				{},
			},
		},
	}

	for streamResponse := range responseChan {
		chunk := []byte(streamResponse.Data)
		response.Data.Choices[0].Content += streamResponse.Data

		if payload.StreamingFunc != nil {
			err := payload.StreamingFunc(ctx, chunk)
			if err != nil {
				return nil, fmt.Errorf("streaming func returned an error: %w", err)
			}
		}
	}
	return &response, nil
}

// 数据格式
// id: "fb981fde-0080-4933-b87b-4a29eaba8d17"
// event: "add"
// data: "作为一个"
//
// id: "fb981fde-0080-4933-b87b-4a29eaba8d17"
// event: "add"
// data: "大型语言模型"
//
// id: "fb981fde-0080-4933-b87b-4a29eaba8d17"
// event: "add"
// data: "我可以"
//
// ... ...
//
// Id: "fb981fde-0080-4933-b87b-4a29eaba8d17"
// event: "finish"
// meta: {"request_id":"123445676789","task_id":"75931252186628","task_status":"SUCCESS","usage":{"prompt_tokens":215,"completion_tokens":302,"total_tokens":517}}
func decodeStreamData(line string) (*StreamedChatResponsePayload, error) {
	if line == "" {
		return nil, errors.New("结果为空")
	}
	re := regexp.MustCompile(`id:\s*"(.*?)"\s*event:\s*"(.*?)"\s*`)
	reData := regexp.MustCompile(`data:\s*"(.*?)"`)
	regMeta := regexp.MustCompile(`meta:\s*({.*?}})`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, errors.New("数据解析失败:" + line)
	}
	resp := &StreamedChatResponsePayload{}
	if matches[2] == eventAdd {
		resp.ID = matches[1]
		resp.Event = matches[2]
		matches1 := reData.FindStringSubmatch(line)
		if len(matches1) != 2 {
			return nil, errors.New("未匹配到data:" + line)
		}
		resp.Data = matches1[1]

	} else if matches[2] == eventError || matches[2] == eventnterrupted {
		return nil, errors.New("error:" + line)
	} else if matches[2] == eventFinish {
		resp.ID = matches[1]
		resp.Event = matches[2]
		matches2 := regMeta.FindStringSubmatch(line)
		if len(matches2) != 2 {
			return nil, errors.New("未匹配到meta:" + line)
		}
		meta := Meta{}
		err := json.Unmarshal([]byte(matches2[1]), &meta)
		if err != nil {
			return nil, errors.New("meta解码失败：" + line)
		}
		resp.Meta = meta
	} else {
		return nil, errors.New("未知的消息事件：" + line)
	}
	return resp, nil
}
