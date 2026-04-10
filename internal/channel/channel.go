package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

// Channel 代表一个飞书机器人实例
type Channel struct {
	AppID     string
	AppSecret string
	Client    *lark.Client
	wsClient  *larkws.Client
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// ChannelManager 管理多个飞书机器人
type ChannelManager struct {
	mu       sync.RWMutex
	Channels map[string]*Channel // key = AppID
}

// NewChannelManager 创建管理器
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		Channels: make(map[string]*Channel),
	}
}

// AddChannel 添加并启动一个机器人（如果已存在则停止并替换）
func (m *ChannelManager) AddChannel(ctx context.Context, appID, appSecret string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果已存在，先停止旧的
	if old, ok := m.Channels[appID]; ok {
		old.Stop()
	}

	channel := &Channel{
		AppID:     appID,
		AppSecret: appSecret,
		Client:    lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug)),
	}

	// 创建事件处理器（闭包捕获 channel，以便回复消息时使用自己的 client）
	handler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			return channel.handleMessage(ctx, event)
		})

	// 创建 WebSocket 客户端
	wsClient := larkws.NewClient(appID, appSecret, larkws.WithEventHandler(handler))
	channel.wsClient = wsClient

	// 启动长连接（独立 goroutine）
	channelCtx, cancel := context.WithCancel(ctx)
	channel.cancel = cancel
	go func() {
		log.Printf("[Channel %s] 启动长连接...", appID)
		if err := wsClient.Start(channelCtx); err != nil {
			log.Printf("[Channel %s] 长连接退出: %v", appID, err)
		}
	}()

	m.Channels[appID] = channel
	log.Printf("[Manager] 机器人 %s 已添加并启动", appID)
	return nil
}

// RemoveChannel 停止并移除机器人
func (m *ChannelManager) RemoveChannel(appID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if channel, ok := m.Channels[appID]; ok {
		channel.Stop()
		delete(m.Channels, appID)
		log.Printf("[Manager] 机器人 %s 已移除", appID)
	}
}

// GetChannel 获取机器人实例
func (m *ChannelManager) GetChannel(appID string) (*Channel, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	channel, ok := m.Channels[appID]
	return channel, ok
}

// StopAll 停止所有机器人
func (m *ChannelManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, channel := range m.Channels {
		channel.Stop()
		delete(m.Channels, id)
	}
}

// Stop 停止当前机器人
func (c *Channel) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}
	log.Printf("[Channel %s] 已停止", c.AppID)
}

// handleMessage 处理收到的消息（与原有逻辑一致，但使用自己的 client 回复）
func (c *Channel) handleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	if event.Event == nil || event.Event.Message == nil {
		log.Printf("[Channel %s] 收到空消息事件", c.AppID)
		return nil
	}
	msg := event.Event.Message
	log.Printf("[Channel %s] 收到消息，chat_id: %s", c.AppID, *msg.ChatId)

	var content struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*msg.Content), &content); err != nil {
		log.Printf("[Channel %s] 解析消息内容失败: %v", c.AppID, err)
		return nil
	}
	log.Printf("[Channel %s] 用户说: %s", c.AppID, content.Text)

	// 回复消息（示例：回显并注明来自哪个机器人）
	replyText := fmt.Sprintf("收到：%s —— 来自机器人 %s", content.Text, c.AppID)
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

	resp, err := c.Client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		log.Printf("[Channel %s] 发送消息 API 调用失败: %v", c.AppID, err)
		return err
	}
	if !resp.Success() {
		log.Printf("[Channel %s] 发送消息业务错误: %s", c.AppID, resp.Msg)
		return fmt.Errorf("send msg error: %s", resp.Msg)
	}
	log.Printf("[Channel %s] 回复成功", c.AppID)
	return nil
}
