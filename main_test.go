package main_test

import (
	"context"
	"github.com/adamchenEpm/ym3-go/internal/config"
	"github.com/adamchenEpm/ym3-go/internal/llm"
	"github.com/adamchenEpm/ym3-go/internal/mysql"
	"github.com/adamchenEpm/ym3-go/internal/redis"
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

/*
 * 测试 mysql.QueryToStructs
 */
func Test_mysql_QueryToStructs(t *testing.T) {
	db := mysql.GetInstance()
	defer db.Close()

	// 4. 查询并扫描到结构体
	type User struct {
		ID         int64      `db:"id"`
		Name       string     `db:"name"`
		UpdateTime *time.Time `db:"update_time"`
	}
	var users []User
	err := db.QueryToStructs("SELECT id, name,update_time FROM sys_user WHERE id = ?", &users, 138)
	if err != nil {
		t.Fatalf("QueryToStructs失败: %v", err)
	}
	if len(users) < 1 {
		t.Logf("查询结果是空的 ")
	} else {
		t.Logf("查询结果正确: %+v", users)
		if users[0].UpdateTime != nil {
			t.Logf("格式化 update_time: %v", users[0].UpdateTime.Format("2006-01-02 15:04:05"))
		}
	}

}

// TestIntegration_BatchInsert 测试批量插入
func TestIntegration_BatchInsert(t *testing.T) {
	db := mysql.GetInstance()
	defer db.Close()

	rows := [][]interface{}{
		{"John", 22},
		{"Jane", 27},
	}
	err := db.BatchInsert("users", []string{"name", "age"}, rows)
	if err != nil {
		t.Fatalf("BatchInsert失败: %v", err)
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
		{Role: "user", Content: "你好，请简单介绍一下自己。"},
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
