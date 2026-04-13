package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/adamchenEpm/ym3-go/internal/channel"
)

func main1() {
	manager := channel.NewChannelManager()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 机器人列表（两个飞书 App）
	// 注意：BotOpenID 需要从飞书开发者后台获取（应用凭证 -> 应用身份 -> open_id）
	botList := []struct {
		AppID     string
		AppSecret string
		BotOpenID string // 机器人的 open_id，用于群聊 @ 检测
	}{
		{
			AppID:     "cli_a945e55a497adbdb",
			AppSecret: "uO2EOunoK5jEQ42ryqfuVbeJGLWe2Prf",
			BotOpenID: "ou_xxxxxx1", // 请替换为实际值
		},
		{
			AppID:     "cli_a95322fcfa391cb6",
			AppSecret: "MyoaM4OKZA8FYK833NmrXggb3X5kFFex",
			BotOpenID: "ou_xxxxxx2", // 请替换为实际值
		},
	}

	for _, cfg := range botList {
		if err := manager.AddChannel(ctx, cfg.AppID, cfg.AppSecret); err != nil {
			log.Printf("添加机器人 %s 失败: %v", cfg.AppID, err)
			continue
		}
		if ch, ok := manager.GetChannel(cfg.AppID); ok {
			ch.SetBotOpenID(cfg.BotOpenID)
			log.Printf("机器人 %s 已设置 BotOpenID = %s", cfg.AppID, cfg.BotOpenID)
			// 可选：关闭“仅@响应”模式（私聊和群聊都会直接回复）
			// ch.OnlyReplyWhenAt = false
		}
	}

	log.Println("所有机器人已启动，等待消息...")
	<-ctx.Done()
	log.Println("关闭所有机器人...")
	manager.StopAll()
}
