package llm

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// 请求结构
// ============================================================================

// ChatCompletionRequest 聊天补全请求（OpenAI 兼容）
type ChatCompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []Message       `json:"messages"`
	Tools            []Tool          `json:"tools,omitempty"`
	ToolChoice       interface{}     `json:"tool_choice,omitempty"` // "none", "auto", "required" 或具体工具
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int  `json:"logit_bias,omitempty"`
	User             string          `json:"user,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	ResponseFormat   *ResponseFormat `json:"response_format,omitempty"`
	Seed             *int            `json:"seed,omitempty"`
}

// ResponseFormat 响应格式
type ResponseFormat struct {
	Type string `json:"type"` // "text" 或 "json_object"
}

// Message 消息结构（支持普通消息、工具调用、工具响应）
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content,omitempty"`      // 消息内容
	Name       string     `json:"name,omitempty"`         // 可选，工具名称或用户名
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // assistant 中的工具调用
	ToolCallID string     `json:"tool_call_id,omitempty"` // tool 响应时关联的调用ID
}

// Validate 校验消息合法性
func (m *Message) Validate() error {
	switch m.Role {
	case "assistant":
		if m.ToolCallID != "" {
			return fmt.Errorf("assistant message cannot have tool_call_id")
		}
		if m.Content == "" && len(m.ToolCalls) == 0 {
			return fmt.Errorf("assistant message must have content or tool_calls")
		}
	case "tool":
		if m.ToolCallID == "" {
			return fmt.Errorf("tool message must have tool_call_id")
		}
		if len(m.ToolCalls) > 0 {
			return fmt.Errorf("tool message cannot have tool_calls")
		}
	case "system", "user":
		if m.ToolCallID != "" || len(m.ToolCalls) > 0 {
			return fmt.Errorf("%s message cannot have tool_calls or tool_call_id", m.Role)
		}
	default:
		return fmt.Errorf("unknown role: %s", m.Role)
	}
	return nil
}

// Tool 工具定义
type Tool struct {
	Type     string   `json:"type"` // 固定 "function"
	Function Function `json:"function"`
}

// Function 函数定义
type Function struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`       // JSON Schema
	Strict      *bool           `json:"strict,omitempty"` // 部分模型支持
}

// ToolCall 工具调用（非流式）
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // "function"
	Function FunctionCall `json:"function"`
}

// FunctionCall 函数调用详情
type FunctionCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"` // JSON 对象字符串
}

// ============================================================================
// 响应结构（非流式）
// ============================================================================

// ChatCompletionResponse 标准聊天补全响应
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	ServiceTier       string   `json:"service_tier,omitempty"`
}

// Choice 选择项
type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	FinishReason string    `json:"finish_reason"` // stop, length, tool_calls, content_filter
	Logprobs     *Logprobs `json:"logprobs,omitempty"`
}

// Logprobs 对数概率信息
type Logprobs struct {
	Content []TokenLogprob `json:"content,omitempty"`
}

// TokenLogprob Token 概率
type TokenLogprob struct {
	Token       string  `json:"token"`
	Logprob     float64 `json:"logprob"`
	Bytes       []int   `json:"bytes,omitempty"`
	TopLogprobs []struct {
		Token   string  `json:"token"`
		Logprob float64 `json:"logprob"`
		Bytes   []int   `json:"bytes,omitempty"`
	} `json:"top_logprobs,omitempty"`
}

// Usage Token 使用统计
type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens,omitempty"`
}

type CompletionTokensDetails struct {
	ReasoningTokens          int `json:"reasoning_tokens,omitempty"`
	AcceptedPredictionTokens int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens int `json:"rejected_prediction_tokens,omitempty"`
}

// ============================================================================
// 流式响应结构
// ============================================================================

// ChatCompletionStreamResponse 流式聊天补全响应（每个 chunk）
type ChatCompletionStreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Usage   *Usage         `json:"usage,omitempty"` // 只在最后一个 chunk 中出现
}

// StreamChoice 流式选择项
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

// StreamDelta 增量消息
type StreamDelta struct {
	Role      string           `json:"role,omitempty"`
	Content   string           `json:"content,omitempty"`
	ToolCalls []ToolCallStream `json:"tool_calls,omitempty"`
}

// ToolCallStream 流式工具调用片段（支持增量拼接）
type ToolCallStream struct {
	Index    int                `json:"index"`
	ID       string             `json:"id,omitempty"`
	Type     string             `json:"type,omitempty"`
	Function FunctionCallStream `json:"function,omitempty"`
}

// FunctionCallStream 流式函数调用（arguments 是字符串片段）
type FunctionCallStream struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"` // 需要累积拼接
}

// ============================================================================
// 错误结构（符合 Go 最佳实践）
// ============================================================================

// APIError OpenAI 兼容的 API 错误
type APIError struct {
	StatusCode int    `json:"-"`               // HTTP 状态码
	Message    string `json:"message"`         // 错误描述
	Type       string `json:"type"`            // 错误类型，如 "invalid_request_error"
	Param      string `json:"param,omitempty"` // 导致错误的参数名
	Code       string `json:"code,omitempty"`  // 错误码，如 "context_length_exceeded"
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	return fmt.Sprintf("llm API error (status=%d): type=%s, code=%s, message=%s",
		e.StatusCode, e.Type, e.Code, e.Message)
}

// Unwrap 支持错误链（可与 errors.Is 配合）
func (e *APIError) Unwrap() error {
	// 这里可以包装原始网络错误等，简化起见返回 nil
	return nil
}

// Is 支持 errors.Is 判断
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	return e.Code == t.Code && e.Type == t.Type
}

// UnmarshalJSON 自定义解析，处理不同后端返回格式
func (e *APIError) UnmarshalJSON(data []byte) error {
	// 尝试标准 OpenAI 格式
	var std struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Param   string `json:"param"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(data, &std); err == nil && std.Error.Message != "" {
		e.Message = std.Error.Message
		e.Type = std.Error.Type
		e.Param = std.Error.Param
		e.Code = std.Error.Code
		return nil
	}

	// 某些代理直接返回 {"error": "some string"}
	var simple struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(data, &simple); err == nil && simple.Error != "" {
		e.Message = simple.Error
		e.Type = "api_error"
		return nil
	}

	// 完全无法解析，保留原始文本
	e.Message = string(data)
	e.Type = "parse_error"
	return nil
}

// ============================================================================
// 辅助函数：构建工具定义
// ============================================================================

// NewTool 便捷创建 Tool
func NewTool(name, description string, parameters interface{}) (Tool, error) {
	var raw json.RawMessage
	switch p := parameters.(type) {
	case string:
		raw = json.RawMessage(p)
	case []byte:
		raw = json.RawMessage(p)
	default:
		b, err := json.Marshal(p)
		if err != nil {
			return Tool{}, fmt.Errorf("marshal parameters: %w", err)
		}
		raw = b
	}
	return Tool{
		Type: "function",
		Function: Function{
			Name:        name,
			Description: description,
			Parameters:  raw,
		},
	}, nil
}

// MustNewTool 创建工具，失败时 panic
func MustNewTool(name, description string, parameters interface{}) Tool {
	t, err := NewTool(name, description, parameters)
	if err != nil {
		panic(err)
	}
	return t
}

// JSONSchemaProperty 定义 JSON Schema 属性
type JSONSchemaProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// NewParameterSchema 构建参数 JSON Schema，返回 json.RawMessage
// 如果发生错误，返回 nil 并记录（调用者需检查）
func NewParameterSchema(properties map[string]JSONSchemaProperty, required []string) (json.RawMessage, error) {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("marshal schema: %w", err)
	}
	return json.RawMessage(data), nil
}

// MustNewParameterSchema 同上，失败 panic
func MustNewParameterSchema(properties map[string]JSONSchemaProperty, required []string) json.RawMessage {
	schema, err := NewParameterSchema(properties, required)
	if err != nil {
		panic(err)
	}
	return schema
}

// ============================================================================
// 辅助函数：响应解析
// ============================================================================

// HasToolCalls 判断响应是否包含工具调用
func (r *ChatCompletionResponse) HasToolCalls() bool {
	if len(r.Choices) == 0 {
		return false
	}
	return len(r.Choices[0].Message.ToolCalls) > 0
}

// GetToolCalls 获取第一个 Choice 的工具调用列表
func (r *ChatCompletionResponse) GetToolCalls() []ToolCall {
	if len(r.Choices) == 0 {
		return nil
	}
	return r.Choices[0].Message.ToolCalls
}

// GetContent 获取第一个 Choice 的消息内容
func (r *ChatCompletionResponse) GetContent() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[0].Message.Content
}

// GetFinishReason 获取第一个 Choice 的结束原因
func (r *ChatCompletionResponse) GetFinishReason() string {
	if len(r.Choices) == 0 {
		return ""
	}
	return r.Choices[0].FinishReason
}

// ============================================================================
// 流式响应辅助：累积工具调用
// ============================================================================

// ToolCallAccumulator 用于累积流式 tool_calls 片段
type ToolCallAccumulator struct {
	calls map[int]*ToolCall // index -> 累积中的 ToolCall
}

// NewToolCallAccumulator 创建累积器
func NewToolCallAccumulator() *ToolCallAccumulator {
	return &ToolCallAccumulator{
		calls: make(map[int]*ToolCall),
	}
}

// AddDelta 添加一个流式 delta 中的 ToolCalls 片段
func (acc *ToolCallAccumulator) AddDelta(deltas []ToolCallStream) {
	for _, delta := range deltas {
		idx := delta.Index
		call, ok := acc.calls[idx]
		if !ok {
			call = &ToolCall{
				ID:   delta.ID,
				Type: delta.Type,
				Function: FunctionCall{
					Name:      delta.Function.Name,
					Arguments: json.RawMessage{}, // 暂空
				},
			}
			acc.calls[idx] = call
		}
		// 更新 ID、Type、Name（可能后续片段才出现）
		if delta.ID != "" {
			call.ID = delta.ID
		}
		if delta.Type != "" {
			call.Type = delta.Type
		}
		if delta.Function.Name != "" {
			call.Function.Name = delta.Function.Name
		}
		// Arguments 是字符串片段，需要累积拼接
		if delta.Function.Arguments != "" {
			// 当前 call.Function.Arguments 是 json.RawMessage，需要转为字符串拼接再转回
			var current string
			if len(call.Function.Arguments) > 0 {
				// 去除可能的双引号包裹（OpenAI 流式返回的 arguments 是裸 JSON 字符串）
				current = string(call.Function.Arguments)
			}
			current += delta.Function.Arguments
			call.Function.Arguments = json.RawMessage(current)
		}
	}
}

// GetCompletedCalls 返回所有已完整累积的 ToolCall（当 FinishReason 非空时调用）
func (acc *ToolCallAccumulator) GetCompletedCalls() []ToolCall {
	result := make([]ToolCall, 0, len(acc.calls))
	for _, call := range acc.calls {
		// 验证必要字段
		if call.ID != "" && call.Function.Name != "" {
			result = append(result, *call)
		}
	}
	return result
}

// Reset 重置累积器
func (acc *ToolCallAccumulator) Reset() {
	acc.calls = make(map[int]*ToolCall)
}

// ============================================================================
// 使用示例（注释）
// ============================================================================

/*
func main() {
    // 1. 定义参数 Schema
    props := map[string]JSONSchemaProperty{
        "location": {Type: "string", Description: "城市名称"},
    }
    schema, _ := NewParameterSchema(props, []string{"location"})

    // 2. 创建工具
    tool, _ := NewTool("get_weather", "查询天气", schema)

    // 3. 构建请求
    req := ChatCompletionRequest{
        Model: "gpt-4o",
        Messages: []Message{
            {Role: "system", Content: "You are a helpful assistant."},
            {Role: "user", Content: "北京天气如何？"},
        },
        Tools:      []Tool{tool},
        ToolChoice: "auto",
    }

    // 4. 发送请求（略）
    // 5. 解析响应
    var resp ChatCompletionResponse
    // ...
    if resp.HasToolCalls() {
        for _, tc := range resp.GetToolCalls() {
            fmt.Println("调用工具:", tc.Function.Name)
            fmt.Println("参数:", string(tc.Function.Arguments))
        }
    } else {
        fmt.Println("回答:", resp.GetContent())
    }
}
*/
