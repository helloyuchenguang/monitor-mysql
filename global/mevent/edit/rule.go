package edit

import (
	"encoding/json"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
	"monitormysql/mgrpc/api/mycanal"
)

type ChannelReplyType interface {
	[]byte | *mycanal.EventTableRowReply
}

type ChannelClient[T ChannelReplyType] struct {
	ID   string
	Chan chan T
}

type RuleServer[T ChannelReplyType] struct {
	clients cmap.ConcurrentMap[string, *ChannelClient[T]]
}

// NewServer  监控处理器接口
func NewServer[T ChannelReplyType]() *RuleServer[T] {
	return &RuleServer[T]{
		clients: cmap.New[*ChannelClient[T]](),
	}
}

// PutNewClient 创建一个新的SSE客户端
func (rs *RuleServer[T]) PutNewClient() *ChannelClient[T] {
	cc := &ChannelClient[T]{
		// 使用 github.com/google/uuid
		ID: uuid.New().String(),
		// 带缓冲防止阻塞
		Chan: make(chan T, 100_000),
	}
	rs.clients.Set(cc.ID, cc)
	slog.Info("添加新客户端", slog.String("clientID", cc.ID))
	return cc
}

func (rs *RuleServer[T]) RemoveClientByID(clientID string) {
	if client, ok := rs.clients.Get(clientID); ok {
		close(client.Chan) // 关闭通道
		rs.clients.Remove(clientID)
		slog.Info("删除客户端", slog.String("clientID", clientID))
	} else {
		slog.Warn("尝试删除不存在的客户端", slog.String("clientID", clientID))
	}
}

func (rs *RuleServer[T]) OnChange(sd *SourceData) error {
	// 发布更新信息到所有 SSE 客户端
	rs.Broadcast(rs.convert(sd))
	return nil
}

func (rs *RuleServer[T]) convert(sd *SourceData) T {
	var data T
	if _, ok := any(data).(*mycanal.EventTableRowReply); ok {
		return any(SourceDataToGrpcReply(sd)).(T)
	}

	if _, ok := any(data).([]byte); ok {
		// 否则使用 []byte 类型
		bs, err := json.Marshal(SourceDataToEditInfo(sd))
		if err != nil {
			slog.Error("JSON 编码错误", slog.Any("error", err))
			return data // 返回零值
		}
		return any(bs).(T)
	}
	slog.Error("未知数据类型", slog.Any("dataType", data))
	return data // 返回零值
}

// Broadcast 广播消息给所有客户端
func (rs *RuleServer[T]) Broadcast(data T) {
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
