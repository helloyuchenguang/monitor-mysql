package monitormysql

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

type SSEClient struct {
	ID   string
	Chan chan []byte
}

type SSEServer struct {
	mu      sync.RWMutex
	clients map[string]*SSEClient
}

func NewSSEServer() *SSEServer {
	return &SSEServer{
		clients: make(map[string]*SSEClient),
	}
}

func (s *SSEServer) AddClient(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	clientID := uuid.New().String() // 使用 github.com/google/uuid
	ch := make(chan []byte, 10)     // 带缓冲防止阻塞

	s.mu.Lock()
	s.clients[clientID] = &SSEClient{ID: clientID, Chan: ch}
	s.mu.Unlock()

	defer func() {
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
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done(): // 客户端断开
			return
		}
	}
}

func (s *SSEServer) Broadcast(data any) {
	msg, _ := json.Marshal(data)

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, client := range s.clients {
		select {
		case client.Chan <- msg:
		default:
			// 防止阻塞：可选择丢弃消息或断开慢客户端
		}
	}
}
