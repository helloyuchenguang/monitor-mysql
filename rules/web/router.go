package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// StartServer 启动Web服务器
func (w *SSERuleService) StartServer() {
	//gin.SetMode(gin.ReleaseMode)
	// gin web
	router := gin.Default()
	router.Static("/static", "./static")

	// 设置index.html为默认页面
	router.GET("/", func(c *gin.Context) {
		c.File("./rules/web/static/index.html")
	})

	// SSE路由
	router.GET("/sse", func(c *gin.Context) {
		rs := w.Rule

		writer := c.Writer
		flusher, ok := writer.(http.Flusher)
		if !ok {
			http.Error(writer, "转换失败", http.StatusInternalServerError)
			return
		}

		// 创建一个新的客户端
		newClient := rs.PutNewClient()
		clientID := newClient.ID
		ch := newClient.Chan
		// 关闭通道
		defer rs.RemoveClientByID(clientID)

		// SSE headers
		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")

		req := c.Request
		for {
			select {
			case dl := <-ch:
				// 发送消息到客户端
				_, err := fmt.Fprintf(writer, "data: %s\n\n", dl)
				if err != nil {
					return
				}
				flusher.Flush()
			case <-req.Context().Done(): // 客户端断开
				slog.Warn("客户端关闭", slog.String("clientID", clientID))
				return
			}
		}
	})
	slog.Info("页面地址: http://localhost" + w.cfg.Addr)
	err := router.Run(w.cfg.Addr)
	if err != nil {
		return
	}
}
