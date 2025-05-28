package edit

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	cmap "github.com/orcaman/concurrent-map/v2"
	"log/slog"
	"monitormysql/mrpc/api/mycanal"
	"net/http"
)

type ChannelReplyType interface {
	[]byte | mycanal.EventTableRowReply
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

// NewClient 创建一个新的SSE客户端
func NewClient[T ChannelReplyType]() *ChannelClient[T] {
	return &ChannelClient[T]{
		// 使用 github.com/google/uuid
		ID: uuid.New().String(),
		// 带缓冲防止阻塞
		Chan: make(chan T, 100),
	}
}

func (rs *RuleServer[T]) OnChange(sd *EditSourceData) error {
	// 发布更新信息到所有 SSE 客户端
	rs.Broadcast(rs.convert(sd))
	return nil
}

func (rs *RuleServer[T]) convert(sd *EditSourceData) T {
	var data T
	if _, ok := any(data).(mycanal.EventTableRowReply); ok {
		panic("暂未实现 mycanal.EventTableRowReply 的转换")
	} else {
		// 否则使用 []byte 类型
		bs, err := json.Marshal(FromRows(sd))
		if err != nil {
			slog.Error("JSON 编码错误", slog.Any("error", err))
			return data // 返回零值
		}
		data = any(bs).(T)
	}
	return data
}

// Broadcast 广播消息给所有客户端
func (rs *RuleServer[T]) Broadcast(data T) {
	for _, client := range rs.clients.Items() {
		select {
		case client.Chan <- data:
		default:
			// 防止阻塞：可选择丢弃消息或断开慢客户端
			slog.Warn("丢弃消息", slog.String("clientID", client.ID))
		}
	}
}

// AddSSEClient 添加客户端
func (rs *RuleServer[T]) AddSSEClient(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "转换失败", http.StatusInternalServerError)
		return
	}

	// 创建一个新的客户端
	newClient := NewClient[T]()
	clientID := newClient.ID
	// 存储客户端到服务器
	rs.clients.Set(clientID, newClient)
	slog.Info("新客户端连接", slog.String("clientID", clientID))

	// 关闭通道
	defer func() {
		// 获取锁，删除客户端
		rs.clients.Remove(clientID)
	}()

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-newClient.Chan:
			// 发送消息到客户端
			_, err := fmt.Fprintf(w, "data: %s\n\n", msg)
			if err != nil {
				return
			}
			flusher.Flush()
		case <-r.Context().Done(): // 客户端断开
			slog.Warn("客户端关闭", slog.String("clientID", clientID))
			return
		}
	}
}
