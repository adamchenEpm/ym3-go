package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func main() {
	// 替换为你的 App ID 和 Secret
	appID := "cli_a945e55a497adbdb"
	appSecret := "uO2EOunoK5jEQ42ryqfuVbeJGLWe2Prf"

	// 1. 初始化 API 客户端（用于发送消息）
	client := lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug), // 可选，便于调试
	)

	// 2. 定义事件处理器（注意：使用 OnP2MessageReceiveV1，对应 Node.js 的 im.message.receive_v1）
	handler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			// 提取消息内容
			if event.Event == nil || event.Event.Message == nil {
				log.Println("⚠️ 消息事件为空")
				return nil
			}
			msg := event.Event.Message
			log.Printf("📨 收到消息，chat_id: %s", *msg.ChatId)

			// 解析消息文本（飞书消息 content 是 JSON 字符串）
			var content struct {
				Text string `json:"text"`
			}
			if err := json.Unmarshal([]byte(*msg.Content), &content); err != nil {
				log.Printf("❌ 解析消息内容失败: %v", err)
				return nil
			}
			log.Printf("💬 用户说: %s", content.Text)

			// 3. 回复消息（与 Node.js 的 client.im.v1.message.create 一致）
			replyText := fmt.Sprintf("回复：%s\n——来自 Go 长连接机器人", content.Text)
			replyContent := map[string]string{"text": replyText}
			replyBytes, _ := json.Marshal(replyContent)

			req := larkim.NewCreateMessageReqBuilder().
				ReceiveIdType("chat_id").
				Body(larkim.NewCreateMessageReqBodyBuilder().
					ReceiveId(*msg.ChatId).
					MsgType("text").
					Content(string(replyBytes)).
					Build()).
				Build()

			resp, err := client.Im.V1.Message.Create(ctx, req)
			if err != nil {
				log.Printf("❌ 发送消息 API 调用失败: %v", err)
				return err
			}
			if !resp.Success() {
				log.Printf("❌ 发送消息业务错误: %s", resp.Msg)
				return fmt.Errorf("send msg error: %s", resp.Msg)
			}
			log.Println("✅ 回复成功")
			return nil
		})

	// 4. 初始化 WebSocket 客户端并启动长连接
	wsClient := larkws.NewClient(appID, appSecret, larkws.WithEventHandler(handler))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("🚀 飞书长连接正在启动...")
	if err := wsClient.Start(ctx); err != nil {
		log.Fatalf("❌ 长连接启动失败: %v", err)
	}
	log.Println("👋 连接已关闭")
}
