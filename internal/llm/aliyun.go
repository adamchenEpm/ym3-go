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

type aliyunClient struct {
	apiKey  string
	baseURL string
	model   string
	httpCli *http.Client
}

// NewAliyunClient 创建阿里百炼客户端
func NewAliyunClient(cfg *Config) Client {
	baseURL := cfg.AliyunBaseURL
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
	model := cfg.AliyunModel
	if model == "" {
		model = "qwen-turbo"
	}
	return &aliyunClient{
		apiKey:  cfg.AliyunAPIKey,
		baseURL: baseURL,
		model:   model,
		httpCli: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *aliyunClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	url := c.baseURL + "/chat/completions"
	model := req.Model
	if model == "" {
		model = c.model
	}

	var tools = []interface{}{
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "get_location",
				"description": "获取用户当前地理位置",
				"parameters": map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
				},
			},
		},
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "get_weather",
				"description": "查询指定城市的天气",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "城市名称，如：北京",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	body := map[string]interface{}{
		"model":       model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      false,
		"tools":       tools,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpCli.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("aliyun error %d: %s", resp.StatusCode, string(bodyBytes))
	} else {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("-----------------")
		fmt.Println(string(bodyBytes))
		fmt.Println("-----------------")
	}

	var apiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&apiResp); err != nil {
		return nil, err
	}
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &ChatResponse{
		Content: apiResp.Choices[0].Message.Content,
		Usage: struct {
			PromptTokens     int
			CompletionTokens int
			TotalTokens      int
		}{
			PromptTokens:     apiResp.Usage.PromptTokens,
			CompletionTokens: apiResp.Usage.CompletionTokens,
			TotalTokens:      apiResp.Usage.TotalTokens,
		},
	}, nil
}
