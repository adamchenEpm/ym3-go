package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
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

	OnlyReplyWhenAt bool   // 群聊时是否仅回复@本机器人的消息
	BotOpenID       string // 机器人的 open_id（用于判断是否被@），需要提前获取
}

// ChannelManager 管理多个飞书机器人
type ChannelManager struct {
	mu       sync.RWMutex
	Channels map[string]*Channel
}

// NewChannelManager 创建管理器
func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		Channels: make(map[string]*Channel),
	}
}

// AddChannel 添加并启动一个机器人
func (m *ChannelManager) AddChannel(ctx context.Context, appID, appSecret string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if old, ok := m.Channels[appID]; ok {
		old.Stop()
	}

	ch := &Channel{
		AppID:           appID,
		AppSecret:       appSecret,
		Client:          lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug)),
		OnlyReplyWhenAt: true, // 默认群聊需要@才回复
		BotOpenID:       "",
	}

	handler := dispatcher.NewEventDispatcher("", "").
		OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			return ch.handleMessage(ctx, event)
		})

	wsClient := larkws.NewClient(appID, appSecret, larkws.WithEventHandler(handler))
	ch.wsClient = wsClient

	chCtx, cancel := context.WithCancel(ctx)
	ch.cancel = cancel
	go func() {
		log.Printf("[Channel %s] 启动长连接...", appID)
		if err := wsClient.Start(chCtx); err != nil {
			log.Printf("[Channel %s] 长连接退出: %v", appID, err)
		}
	}()

	m.Channels[appID] = ch
	log.Printf("[Manager] 机器人 %s 已添加并启动", appID)
	return nil
}

// RemoveChannel 停止并移除机器人
func (m *ChannelManager) RemoveChannel(appID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ch, ok := m.Channels[appID]; ok {
		ch.Stop()
		delete(m.Channels, appID)
		log.Printf("[Manager] 机器人 %s 已移除", appID)
	}
}

// GetChannel 获取机器人实例
func (m *ChannelManager) GetChannel(appID string) (*Channel, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ch, ok := m.Channels[appID]
	return ch, ok
}

// StopAll 停止所有机器人
func (m *ChannelManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, ch := range m.Channels {
		ch.Stop()
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

// SetBotOpenID 设置机器人的 open_id（用于群聊@检测）
func (c *Channel) SetBotOpenID(openID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.BotOpenID = openID
}

// handleMessage 处理消息
func (c *Channel) handleMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	if event.Event == nil || event.Event.Message == nil {
		log.Printf("[Channel %s] 收到空消息事件", c.AppID)
		return nil
	}
	msg := event.Event.Message

	chatID := *msg.ChatId
	chatType := *msg.ChatType

	var senderID string
	if event.Event.Sender != nil && event.Event.Sender.SenderId != nil {
		sender := event.Event.Sender.SenderId
		if sender.UserId != nil {
			senderID = *sender.UserId
		} else if sender.OpenId != nil {
			senderID = *sender.OpenId
		}
	}

	var content struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(*msg.Content), &content); err != nil {
		log.Printf("[Channel %s] 解析消息内容失败: %v", c.AppID, err)
		return nil
	}
	rawText := content.Text
	log.Printf("[Channel %s] 收到消息 | type=%s | chatID=%s | sender=%s | BotOpenID=%s | text=%s",
		c.AppID, chatType, chatID, senderID, c.BotOpenID, rawText)

	finalText := rawText
	if chatType == "group" && c.OnlyReplyWhenAt {
		atMe, cleaned := c.checkAtMe(rawText)
		if !atMe {
			log.Printf("[Channel %s] 群消息未@本机器人，忽略", c.AppID)
			return nil
		}
		finalText = cleaned
	}

	replyText := fmt.Sprintf("收到：%s —— 来自机器人 %s", finalText, c.AppID)
	replyContent := map[string]string{"text": replyText}
	replyBytes, _ := json.Marshal(replyContent)

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType("chat_id").
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(chatID).
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

// checkAtMe 检查消息中是否 @ 了本机器人，并移除 @ 标签
func (c *Channel) checkAtMe(text string) (bool, string) {
	if c.BotOpenID == "" {
		// 未配置机器人 open_id，默认当作已@（所有消息都响应）
		log.Printf("[Channel %s] BotOpenID 未配置，默认响应所有群消息", c.AppID)
		return true, text
	}

	// 1. 正则提取飞书 @ 标签中的 user_id
	// 格式: <at user_id="ou_xxx">...</at> 或 <at user_id='ou_xxx'>...</at>
	re := regexp.MustCompile(`<at user_id="([^"]+)"[^>]*>.*?</at>`)
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		userID := matches[1]
		if userID == c.BotOpenID {
			// 移除整个 <at> 标签及其内容，得到纯净文本
			cleaned := re.ReplaceAllString(text, "")
			cleaned = strings.TrimSpace(cleaned)
			log.Printf("[Channel %s] 检测到 XML @ 标签，匹配成功，清理后文本: %s", c.AppID, cleaned)
			return true, cleaned
		}
	}

	// 2. 兼容纯文本 @open_id 形式（某些情况下飞书可能直接渲染成文本）
	if strings.Contains(text, "@"+c.BotOpenID) {
		cleaned := strings.Replace(text, "@"+c.BotOpenID, "", 1)
		cleaned = strings.TrimSpace(cleaned)
		log.Printf("[Channel %s] 检测到纯文本 @ 匹配，清理后文本: %s", c.AppID, cleaned)
		return true, cleaned
	}

	log.Printf("[Channel %s] 未检测到 @ 本机器人，BotOpenID=%s", c.AppID, c.BotOpenID)
	return false, text
}
