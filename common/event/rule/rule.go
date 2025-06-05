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

type Server struct {
	clients map[string]*ChannelClient
}

// NewServer  监控处理器接口
func NewServer() *Server {
	return &Server{
		clients: make(map[string]*ChannelClient),
	}
}

// PutNewClient 创建一个新的SSE客户端
func (rs *Server) PutNewClient() *ChannelClient {
	cc := &ChannelClient{
		// 使用 github.com/google/uuid
		ID: uuid.New().String(),
		// 带缓冲防止阻塞
		Chan: make(chan *event.Data, 1_000),
	}
	rs.clients[cc.ID] = cc
	slog.Info("添加新客户端", slog.String("clientID", cc.ID))
	return cc
}

func (rs *Server) RemoveClientByID(clientID string) {
	if client, ok := rs.clients[clientID]; ok {
		close(client.Chan) // 关闭通道
		// 从客户端列表中删除
		delete(rs.clients, clientID)
		slog.Info("删除客户端", slog.String("clientID", clientID))
	} else {
		slog.Warn("尝试删除不存在的客户端", slog.String("clientID", clientID))
	}
}

func (rs *Server) OnNext(data *event.Data) error {
	if len(rs.clients) == 1 {
		// 如果只有一个客户端，直接发送数据
		for _, client := range rs.clients {
			forward(client, data)
			return nil
		}
	} else {
		for _, client := range rs.clients {
			go func(c *ChannelClient) {
				forward(c, data)
			}(client)
		}
	}
	return nil
}

// forward 将数据发送到客户端的通道
func forward(client *ChannelClient, data *event.Data) {
	select {
	case client.Chan <- data:
	default:
		// 防止阻塞：可选择丢弃消息或断开慢客户端
		slog.Warn("丢弃消息", slog.String("clientID", client.ID))
	}
}

func (rs *Server) ClientIsEmpty() bool {
	// 检查是否有客户端连接
	return len(rs.clients) == 0
}
