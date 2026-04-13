package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type customClient struct {
	endpoint   string
	apiKey     string
	modelName  string
	headers    map[string]string
	reqMapper  func(interface{}) ([]byte, error)
	respMapper func([]byte) (string, error)
	httpCli    *http.Client
}

// NewCustomClient 创建自定义通用客户端
func NewCustomClient(cfg *Config) Client {
	return &customClient{
		endpoint:   cfg.CustomEndpoint,
		apiKey:     cfg.CustomAPIKey,
		modelName:  cfg.CustomModelName,
		headers:    cfg.CustomHeaders,
		reqMapper:  cfg.CustomReqMapper,
		respMapper: cfg.CustomRespMapper,
		httpCli:    &http.Client{Timeout: 60 * time.Second},
	}
}

// defaultReqMapper 默认请求映射：构造 OpenAI 兼容格式
func defaultReqMapper(req *ChatRequest, modelName string) ([]byte, error) {
	body := map[string]interface{}{
		"model":       modelName,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      false,
	}
	return json.Marshal(body)
}

// defaultRespMapper 默认响应映射：解析 OpenAI 兼容格式
func defaultRespMapper(data []byte) (string, error) {
	var resp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices")
	}
	return resp.Choices[0].Message.Content, nil
}

func (c *customClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 使用自定义请求映射或默认
	var bodyBytes []byte
	var err error
	if c.reqMapper != nil {
		bodyBytes, err = c.reqMapper(req)
	} else {
		bodyBytes, err = defaultReqMapper(req, c.modelName)
	}
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	for k, v := range c.headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("custom api error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var content string
	if c.respMapper != nil {
		content, err = c.respMapper(respData)
	} else {
		content, err = defaultRespMapper(respData)
	}
	if err != nil {
		return nil, err
	}

	// 自定义客户端无法获取 token 使用量，置零
	return &ChatResponse{
		Content: content,
		Usage: struct {
			PromptTokens     int
			CompletionTokens int
			TotalTokens      int
		}{0, 0, 0},
	}, nil
}
