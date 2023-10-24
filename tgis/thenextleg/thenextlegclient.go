package thenextleg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultBaseUrl = "https://api.thenextleg.io/v2"
)

type TheNextLeg struct {
	baseUrl    string // 请求url
	authToken  string
	httpClient Doer
}

type Option func(*TheNextLeg)

func WithAuthToken(token string) Option {
	return func(leg *TheNextLeg) {
		leg.authToken = token
	}
}

// Doer performs a HTTP request.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// New 创建TheNextLeg 客户端实例
func New(opts ...Option) (*TheNextLeg, error) {
	c := &TheNextLeg{}
	for _, v := range opts {
		v(c)
	}
	if c.baseUrl == "" {
		c.baseUrl = defaultBaseUrl
	}
	if c.authToken == "" {
		return nil, errors.New("缺少token")
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c, nil
}

type ImagineRequest struct {
	Msg             string `json:"msg"`
	Ref             string `json:"ref,omitempty"`
	WebhookOverride string `json:"webhookOverride,omitempty"`
	IgnorePrefilter string `json:"ignorePrefilter,omitempty"`
}
type MsgIdResponse struct {
	Success   bool   `json:"success"`
	Msg       string `json:"msg,omitempty"`
	MessageId string `json:"messageId,omitempty"`
	CreateAt  string `json:"createAt,omitempty"`
}

// Imagine 生成图片
func (t *TheNextLeg) Imagine(ctx context.Context, payload *ImagineRequest) (*MsgIdResponse, error) {
	url := fmt.Sprintf("%s/imagine", t.baseUrl)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var resp MsgIdResponse
	if err := t.doHttp(ctx, url, http.MethodPost, payloadBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MessageResponse 生成图片的response
type MessageResponse struct {
	Progress         any      `json:"progress"` // success  100 , in progress : 37, incomplete : incomplete
	Response         Response `json:"response,omitempty"`
	ProgressImageUrl string   `json:"progressImageUrl,omitempty"` // in progress

}
type Response struct {
	CreatedAt            string   `json:"createdAt,omitempty"`
	Buttons              []string `json:"buttons,omitempty"`
	ImageUrl             string   `json:"imageUrl,omitempty"`
	ImageUrls            []string `json:"imageUrls,omitempty"`
	ButtonMessageId      string   `json:"buttonMessageId,omitempty"`
	OriginatingMessageId string   `json:"originatingMessageId,omitempty"`
	Content              string   `json:"content,omitempty"`
	Ref                  string   `json:"ref,omitempty"`
	ResponseAt           string   `json:"responseAt,omitempty"`
}

//// ButtonMessageResponse buttonMessageId的response
//type ButtonMessageResponse struct {
//	CreatedAt struct {
//		Nanoseconds int64 `json:"_nanoseconds"`
//		Seconds     int   `json:"_seconds"`
//	} `json:"createdAt"`
//	Buttons              []string `json:"buttons"`
//	Type                 string   `json:"type"`
//	ImageUrl             string   `json:"imageUrl"`
//	ButtonMessageId      string   `json:"buttonMessageId"`
//	OriginatingMessageId string   `json:"originatingMessageId"`
//	Content              string   `json:"content"`
//	Ref                  string   `json:"ref"`
//	ResponseAt           string   `json:"responseAt"`
//}

// Message 获取任务进入，progress 值100：success， 37：生成进度， "incomplete"：失败，未完成
func (t *TheNextLeg) Message(ctx context.Context, msgId string) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/message/%s", t.baseUrl, msgId)

	var resp MessageResponse
	if err := t.doHttp(ctx, url, http.MethodGet, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// MessageButton 获取任务进入，progress 值100：success， 37：生成进度， "incomplete"：失败，未完成
func (t *TheNextLeg) MessageButton(ctx context.Context, msgId string) (*MessageResponse, error) {
	url := fmt.Sprintf("%s/message/%s", t.baseUrl, msgId)

	var resp MessageResponse
	if err := t.doHttp(ctx, url, http.MethodGet, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type ButtonRequest struct {
	ButtonMessageId string `json:"buttonMessageId"`
	Button          string `json:"button,omitempty"`
	Ref             string `json:"ref,omitempty"`
	WebhookOverride string `json:"webhookOverride,omitempty"`
	Prompt          string `json:"prompt,omitempty"`
	Zoom            string `json:"zoom,omitempty"`
	Ar              string `json:"ar,omitempty"`
}

// Button 根据返回的额button按钮继续操作
func (t *TheNextLeg) Button(ctx context.Context, payload *ButtonRequest) (*MsgIdResponse, error) {
	url := fmt.Sprintf("%s/button", t.baseUrl)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var resp MsgIdResponse
	if err := t.doHttp(ctx, url, http.MethodPost, payloadBytes, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (t *TheNextLeg) doHttp(ctx context.Context, url, method string, payloadBytes []byte, resp any) error {
	// Build request
	var body io.Reader
	if payloadBytes != nil {
		body = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	t.setHeader(req)
	r, err := t.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(r.Body)
	if r.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("API returned unexpected status code: %d", r.StatusCode)
		return errors.New(msg)
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	if err := json.Unmarshal(b, resp); err != nil {
		return err
	}
	return nil
}

func (t *TheNextLeg) setHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+t.authToken)
	req.Header.Set("Content-Type", "application/json")
}
