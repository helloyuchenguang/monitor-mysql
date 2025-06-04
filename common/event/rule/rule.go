package rule

import (
	"github.com/google/uuid"
	"log/slog"
	"main/common/event"
)

// MonitorRuler 监控处理器接口
type MonitorRuler interface {
	// OnNext 处理更新信息
	OnNext(d *event.Data) error

	ClientIsEmpty() bool
}

type ChannelClient struct {
	ID   string
	Chan chan *event.Data
}

type RuleServer struct {
	clients map[string]*ChannelClient
}

// NewServer  监控处理器接口
func NewServer() *RuleServer {
	return &RuleServer{
		clients: make(map[string]*ChannelClient),
	}
}

// PutNewClient 创建一个新的SSE客户端
func (rs *RuleServer) PutNewClient() *ChannelClient {
	cc := &ChannelClient{
		// 使用 github.com/google/uuid
		ID: uuid.New().String(),
		// 带缓冲防止阻塞
		Chan: make(chan *event.Data, 100_000),
	}
	rs.clients[cc.ID] = cc
	slog.Info("添加新客户端", slog.String("clientID", cc.ID))
	return cc
}

func (rs *RuleServer) RemoveClientByID(clientID string) {
	if client, ok := rs.clients[clientID]; ok {
		close(client.Chan) // 关闭通道
		// 从客户端列表中删除
		delete(rs.clients, clientID)
		slog.Info("删除客户端", slog.String("clientID", clientID))
	} else {
		slog.Warn("尝试删除不存在的客户端", slog.String("clientID", clientID))
	}
}

func (rs *RuleServer) OnNext(data *event.Data) error {
	// 发布更新信息到所有客户端
	rs.Broadcast(data)
	return nil
}

func (rs *RuleServer) ClientIsEmpty() bool {
	// 检查是否有客户端连接
	return len(rs.clients) == 0
}

// Broadcast 广播消息给所有客户端
func (rs *RuleServer) Broadcast(data *event.Data) {
	for _, client := range rs.clients {
		go func() {
			select {
			case client.Chan <- data:
			default:
				// 防止阻塞：可选择丢弃消息或断开慢客户端
				slog.Warn("丢弃消息", slog.String("clientID", client.ID))
			}
		}()
	}
}
