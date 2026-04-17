package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adamchenEpm/ym3-go/internal/common"
	"github.com/adamchenEpm/ym3-go/internal/config"
	"github.com/adamchenEpm/ym3-go/internal/llm"
	"github.com/adamchenEpm/ym3-go/internal/pg"
	pgmodel "github.com/adamchenEpm/ym3-go/internal/pg/model"
	"github.com/adamchenEpm/ym3-go/internal/redis"
	"github.com/joho/godotenv"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

/*
 * 测试 config.NewConfig
 */
func Test_config_NewConfig(t *testing.T) {

	cfg := config.Get()
	//t.Assert(cfg != nil)

	t.Logf("Config.name: %v,  code :%v", cfg.Name, cfg.Code)
}

func Test_env(t *testing.T) {

	_ = godotenv.Load()
	apiKey := os.Getenv("LLM_ALIYUN_API_KEY")
	t.Logf("LLM_ALIYUN_API_KEY: %v", apiKey)

}

func Test_common_Encrypt(t *testing.T) {
	pwd := "123"
	encrypt, err := common.Encrypt(pwd)
	if err != nil {
		t.Fatalf("Encrypt失败: %v", err)
	}
	t.Logf("加密后 : %v ", encrypt)

	decrypt, err := common.Decrypt(encrypt)
	t.Logf("解密后 : %v ", decrypt)

}

func Test_exec(t *testing.T) {

	type Request struct {
		City string `json:"city"`
	}

	// 1. 构造要传给 Python 的数据
	req := Request{
		City: "上海", // 你可以改成任意城市：上海、广州...
	}

	// 2. 序列化为 JSON 字符串
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Println("JSON 序列化失败:", err)
		return
	}
	jsonStr := string(jsonBytes)
	fmt.Println("传入参数:", jsonStr)

	// ==================== 核心：跨平台执行 Python 脚本 ====================
	var cmd *exec.Cmd

	// Windows 和 Linux/Mac 命令不一样，自动判断系统
	if runtime.GOOS == "windows" {
		// Windows：python 脚本.py "JSON"
		cmd = exec.Command("python", "e:/app1/ym3-go/data/agents/ym3_weather/weather.py")
		stdin, _ := cmd.StdinPipe()
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, jsonStr)
		}()

		out, _ := cmd.CombinedOutput()

		fmt.Println(string(out))
		fmt.Println("================ 天气结果py ================")

		cmd = exec.Command("node", "e:/app1/ym3-go/data/agents/ym3_weather/weather.js")
		stdin, _ = cmd.StdinPipe()
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, jsonStr)
		}()

		out, _ = cmd.CombinedOutput()

		fmt.Println(string(out))
		fmt.Println("================ 天气结果js ================")

	} else {
		// Linux / Mac：python3 脚本.py "JSON"
		cmd = exec.Command("python3", "weather.py", jsonStr)
	}

}

func Test_pg_QueryToStructs(t *testing.T) {
	db := pg.GetInstance()
	defer func() { _ = db.Close() }()

	var llms []pgmodel.SysLlm

	err := db.QueryToStructs(pgmodel.SysLlmSelect+" where id = $1", &llms, 1)
	if err != nil {
		t.Fatalf("QueryToStructs失败: %v", err)
	}

	if len(llms) == 0 {
		t.Logf("查询结果是空的 ")
	} else {
		t.Logf("查询结果正确: %+v ", llms[0])
		t.Logf("更新时间:%s", common.TimeToStr(llms[0].UpdateTime))
	}

}

func Test_Redis_BasicOps(t *testing.T) {
	rdb := redis.GetInstance()
	defer rdb.Close()

	key := "test:user:138"
	// 1. 设置值
	err := rdb.Set(key, "张三", 60*time.Second)
	if err != nil {
		t.Fatalf("Set失败: %v", err)
	}

	// 2. 获取值
	val, err := rdb.Get(key)
	if err != nil {
		t.Fatalf("Get失败: %v", err)
	}
	if val != "张三" {
		t.Errorf("期望 '张三', 得到 '%s'", val)
	}
	t.Logf("Get成功: %s = %s", key, val)

	// 3. 删除
	//err = rdb.Del(key)
	//if err != nil {
	//	t.Fatalf("Del失败: %v", err)
	//}
	//val2, err := rdb.Get(key)
	//if err == nil {
	//	t.Errorf("期望 key 不存在，但得到 '%s'", val2)
	//}
	//t.Log("删除后 key 已不存在")
}

func Test_llm_Aliyun(t *testing.T) {
	// 1. 创建阿里百炼客户端配置
	cfg := &llm.Config{
		Provider:      llm.ProviderAliyun,
		AliyunAPIKey:  "sk-16ab6965525b4bd4bd245d9e8a3a693c", // 请替换为您的真实 API Key
		AliyunBaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		AliyunModel:   "glm-5", // 可选 qwen-plus, qwen-max
	}

	// 2. 创建客户端
	client, err := llm.NewClient(cfg)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}

	// 3. 构造请求消息
	messages := []llm.Message{
		{Role: "system", Content: "你是一个乐于助人的助手，请用中文回答。"},
		{Role: "user", Content: "你好，你能告诉我广州的天气吗"},
	}

	req := &llm.ChatRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   500,
		// Model 字段可选，如果留空则使用配置中的 AliyunModel
	}

	// 4. 调用大模型（带超时控制）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("调用大模型失败: %v", err)
	}

	// 5. 输出结果
	t.Logf("模型回复: %s", resp.Content)
	t.Logf("Token使用情况: prompt=%d, completion=%d, total=%d",
		resp.Usage.PromptTokens,
		resp.Usage.CompletionTokens,
		resp.Usage.TotalTokens)

	// 可选：验证非空
	if resp.Content == "" {
		t.Errorf("回复内容为空")
	}
}
