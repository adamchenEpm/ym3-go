package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamchenEpm/ym3-go/internal/llm"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// 配置
	_ = godotenv.Load()

	apiKey := os.Getenv("LLM_ALIYUN_API_KEY") // 从环境变量读取
	baseURL := "https://dashscope.aliyuncs.com/compatible-mode/v1"

	// ========== 示例 1：普通对话 + 工具调用 ==========
	fmt.Println("=== 示例 1：普通对话（带工具） ===")
	exampleNormalCall(apiKey, baseURL)

	// ========== 示例 2：流式对话 ==========
	fmt.Println("\n=== 示例 2：流式对话 ===")
	exampleStreamCall(apiKey, baseURL)

}

// 普通调用示例（包含工具调用处理）
func exampleNormalCall(apiKey, baseURL string) {
	// 1. 定义工具参数 Schema
	weatherProps := map[string]llm.JSONSchemaProperty{
		"location": {
			Type:        "string",
			Description: "城市名称，如北京、上海",
		},
		"unit": {
			Type:        "string",
			Description: "温度单位",
			Enum:        []string{"celsius", "fahrenheit"},
		},
	}
	weatherSchema, err := llm.NewParameterSchema(weatherProps, []string{"location"})
	if err != nil {
		panic(err)
	}

	locationSchema, _ := llm.NewParameterSchema(map[string]llm.JSONSchemaProperty{}, nil)

	// 2. 创建工具
	weatherTool, _ := llm.NewTool("get_weather", "查询指定城市的天气", weatherSchema)
	locationTool, _ := llm.NewTool("get_location", "获取用户当前地理位置", locationSchema)

	// 3. 构建请求
	req := llm.ChatCompletionRequest{
		Model: "glm-5",
		Messages: []llm.Message{
			{Role: "system", Content: "你是一个智能助手，可以使用提供的工具帮助用户。"},
			{Role: "user", Content: "我这里天气怎么样？"},
		},
		Tools:       []llm.Tool{weatherTool, locationTool},
		ToolChoice:  "auto",
		Temperature: ptrFloat64(0.7),
	}

	// 4. 发送请求
	resp, err := sendRequest(apiKey, baseURL+"/chat/completions", req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return
	}

	// 5. 处理响应
	if resp.HasToolCalls() {
		fmt.Println("模型要求调用以下工具：")
		for _, tc := range resp.GetToolCalls() {
			fmt.Printf("  - 工具名: %s\n", tc.Function.Name)
			fmt.Printf("    参数: %s\n", string(tc.Function.Arguments))
			// 这里应该执行实际工具并再次调用模型，本例仅演示
		}
	} else {
		fmt.Printf("模型回复: %s\n", resp.GetContent())
	}
	fmt.Printf("Token使用: prompt=%d, completion=%d, total=%d\n",
		resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
}

// 流式调用示例
func exampleStreamCall(apiKey, baseURL string) {
	// 构建请求（不包含工具，仅演示流式）
	req := llm.ChatCompletionRequest{
		Model: "glm-5",
		Messages: []llm.Message{
			{Role: "system", Content: "你是一个助手，用中文回答。"},
			{Role: "user", Content: "请用三句话介绍一下人工智能。"},
		},
		Stream:      true,
		Temperature: ptrFloat64(0.5),
	}

	// 发送流式请求
	err := sendStreamRequest(apiKey, baseURL, req, func(chunk llm.ChatCompletionStreamResponse) {
		if len(chunk.Choices) > 0 {
			delta := chunk.Choices[0].Delta
			fmt.Print(delta.Content) // 实时打印
		}
	})
	if err != nil {
		fmt.Printf("流式请求失败: %v\n", err)
	}
	fmt.Println() // 换行
}

// ========== 底层 HTTP 发送函数 ==========

// sendRequest 发送非流式请求，返回解析后的响应
func sendRequest(apiKey, baseURL string, req llm.ChatCompletionRequest) (*llm.ChatCompletionResponse, error) {
	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer httpResp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		var apiErr llm.APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, fmt.Errorf("status %d: %s", httpResp.StatusCode, string(respBody))
		}
		apiErr.StatusCode = httpResp.StatusCode
		return nil, &apiErr
	}

	// 解析成功响应
	var chatResp llm.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	// 验证消息合法性（可选）
	for _, choice := range chatResp.Choices {
		if err := choice.Message.Validate(); err != nil {
			return nil, fmt.Errorf("invalid message: %w", err)
		}
	}

	return &chatResp, nil
}

// sendStreamRequest 发送流式请求，通过回调处理每个 chunk
func sendStreamRequest(apiKey, baseURL string, req llm.ChatCompletionRequest, onChunk func(llm.ChatCompletionStreamResponse)) error {
	// 强制设置 Stream = true
	req.Stream = true

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{Timeout: 60 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return fmt.Errorf("status %d: %s", httpResp.StatusCode, string(body))
	}

	// 逐行读取 SSE 数据
	scanner := bufio.NewScanner(httpResp.Body)
	// 设置更大的缓冲区（某些响应较长）
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk llm.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			// 忽略解析错误的 chunk，继续
			continue
		}
		onChunk(chunk)
	}
	return scanner.Err()
}

// 辅助函数：返回 float64 指针
func ptrFloat64(v float64) *float64 {
	return &v
}
