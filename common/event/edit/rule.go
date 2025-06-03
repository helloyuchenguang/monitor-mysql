package edit

import (
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
	"main/common/event"
	"main/rules/mgrpc/api/mycanal"
)

type ChannelReplyType interface {
	[]byte | *mycanal.EventTableRowReply
}

type ChannelClient struct {
	ID   string
	Chan chan *event.Data
}

type RuleServer struct {
	clients cmap.ConcurrentMap[string, *ChannelClient]
}

// NewServer  监控处理器接口
func NewServer() *RuleServer {
	return &RuleServer{
		clients: cmap.New[*ChannelClient](),
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
	rs.clients.Set(cc.ID, cc)
	slog.Info("添加新客户端", slog.String("clientID", cc.ID))
	return cc
}

func (rs *RuleServer) RemoveClientByID(clientID string) {
	if client, ok := rs.clients.Get(clientID); ok {
		close(client.Chan) // 关闭通道
		rs.clients.Remove(clientID)
		slog.Info("删除客户端", slog.String("clientID", clientID))
	} else {
		slog.Warn("尝试删除不存在的客户端", slog.String("clientID", clientID))
	}
}

func (rs *RuleServer) OnChange(data *event.Data) error {
	// 发布更新信息到所有客户端
	rs.Broadcast(data)
	return nil
}

// Broadcast 广播消息给所有客户端
func (rs *RuleServer) Broadcast(data *event.Data) {
	for _, client := range rs.clients.Items() {
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
