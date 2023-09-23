package qwenclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	defaultEmbeddingModel = "text-embedding-v1" // 1536 25 2048
)

type EmbeddingPayload struct {
	Model     string   `json:"model"`
	Input     []string `json:"input"`
	TextType  string   `json:"text_type"` //取值：query 或者 document，默认值为 document  ,存库使用document， 查询query
	RequestId string   `json:"request_id,omitempty"`
}

type embeddingResponsePayload struct {
	Output struct {
		Embeddings []*struct {
			TextIndex int       `json:"text_index"`
			Embedding []float64 `json:"embedding"`
		} `json:"embeddings"`
	} `json:"output"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	}
	RequestId string `json:"request_id"`
}

func (c *Client) createEmbedding(ctx context.Context, payload *EmbeddingPayload) (*embeddingResponsePayload, error) {
	fmt.Printf("%#v\n", payload)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	if payload.Model == "" {
		c.embeddingsModel = defaultEmbeddingModel
	}
	if c.baseURL == "" {
		c.baseURL = defaultEmbeddingURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeader(req, false)

	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
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

	var response embeddingResponsePayload

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	fmt.Printf("%#v\n", response)
	return &response, nil
}
