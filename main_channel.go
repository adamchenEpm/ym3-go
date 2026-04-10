package main

import (
	"context"
	"github.com/adamchenEpm/ym3-go/internal/channel"
	"log"
	"os/signal"
	"syscall"
)

// 模拟从数据库或配置文件加载的机器人列表
var botConfigs = []struct {
	AppID     string
	AppSecret string
}{
	{
		AppID:     "cli_a945e55a497adbdb",
		AppSecret: "uO2EOunoK5jEQ42ryqfuVbeJGLWe2Prf",
	},
	{
		AppID:     "cli_a95322fcfa391cb6",
		AppSecret: "MyoaM4OKZA8FYK833NmrXggb3X5kFFex",
	},
}

func main() {
	manager := channel.NewChannelManager()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 注册所有机器人
	for _, cfg := range botConfigs {
		if err := manager.AddChannel(ctx, cfg.AppID, cfg.AppSecret); err != nil {
			log.Printf("添加渠道 %s 失败: %v", cfg.AppID, err)
		}
	}

	log.Println("所有渠道已启动，等待消息...")
	<-ctx.Done()
	log.Println("收到退出信号，关闭所有渠道连接...")
	manager.StopAll()
	log.Println("程序退出")
}
