package llm

import (
	"fmt"
)

// Provider 定义支持的提供商类型
type Provider string

const (
	ProviderOpenAI Provider = "openai" // OpenAI 兼容（DeepSeek, Moonshot 等）
	ProviderAliyun Provider = "aliyun" // 阿里百炼
	ProviderCustom Provider = "custom" // 自定义通用
)

// Config 模型客户端配置
type Config struct {
	Provider Provider `json:"provider"`

	// OpenAI 兼容配置（适用于 OpenAI、DeepSeek、Moonshot 等）
	OpenAICompatBaseURL string `json:"openai_compat_base_url"` // 如 https://api.openai.com/v1
	OpenAICompatAPIKey  string `json:"openai_compat_api_key"`
	OpenAICompatModel   string `json:"openai_compat_model"` // 如 gpt-3.5-turbo, deepseek-chat

	// 阿里百炼配置
	AliyunAPIKey  string `json:"aliyun_api_key"`
	AliyunBaseURL string `json:"aliyun_base_url"` // 默认 https://dashscope.aliyuncs.com/compatible-mode/v1
	AliyunModel   string `json:"aliyun_model"`    // 如 qwen-turbo, qwen-plus

	// 自定义通用配置
	CustomEndpoint   string                            `json:"custom_endpoint"`   // 完整的 API 地址
	CustomAPIKey     string                            `json:"custom_api_key"`    // 可选，会放在 Authorization: Bearer 头
	CustomModelName  string                            `json:"custom_model_name"` // 请求中使用的模型名
	CustomHeaders    map[string]string                 `json:"custom_headers"`    // 额外请求头
	CustomReqMapper  func(interface{}) ([]byte, error) `json:"-"`                 // 自定义请求转换函数
	CustomRespMapper func([]byte) (string, error)      `json:"-"`                 // 自定义响应转换函数
}

// NewClient 根据配置创建对应的模型客户端
func NewClient(cfg *Config) (Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	switch cfg.Provider {
	case ProviderOpenAI:
		return NewOpenAICompatClient(cfg), nil
	case ProviderAliyun:
		return NewAliyunClient(cfg), nil
	case ProviderCustom:
		return NewCustomClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", cfg.Provider)
	}
}
