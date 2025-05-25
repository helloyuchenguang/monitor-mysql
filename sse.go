package monitormysql

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"sync"
)

// SSEClient 代表一个客户端连接
type SSEClient struct {
	ID   string
	Chan chan []byte
}

func NewSSEClient() *SSEClient {
	return &SSEClient{
		// 使用 github.com/google/uuid
		ID: uuid.New().String(),
		// 带缓冲防止阻塞
		Chan: make(chan []byte, 10),
	}
}

// SSEServer 代表一个SSE服务器
type SSEServer struct {
	mu      sync.RWMutex
	clients map[string]*SSEClient
}

// NewSSEServer  监控处理器接口
func NewSSEServer() *SSEServer {
	return &SSEServer{
		clients: make(map[string]*SSEClient),
	}
}

// AddClient 添加客户端
func (s *SSEServer) AddClient(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 创建一个新的客户端
	s.mu.Lock()
	newClient := NewSSEClient()
	clientID := newClient.ID
	s.clients[clientID] = newClient
	s.mu.Unlock()

	// 关闭通道
	defer func() {
		// 获取锁，删除客户端
		s.mu.Lock()
		delete(s.clients, clientID)
		s.mu.Unlock()
	}()

	// SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg := <-newClient.Chan:
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

// Broadcast 广播消息给所有客户端
func (s *SSEServer) Broadcast(data any) {
	msg, _ := json.Marshal(data)

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client := range s.clients {
		select {
		case client.Chan <- msg:
		default:
			// 防止阻塞：可选择丢弃消息或断开慢客户端
			slog.Warn("丢弃消息", slog.String("clientID", client.ID))
		}
	}
}
