package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"monitormysql/global"
	"net/http"
)

// StartServer 启动Web服务器
func StartServer(addr string) {
	// gin web
	router := gin.Default()
	router.Static("/static", "./static")

	// 设置index.html为默认页面
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// SSE路由
	router.GET("/sse", func(c *gin.Context) {
		rs := global.GetRule[[]byte](RuleName)

		writer := c.Writer
		flusher, ok := writer.(http.Flusher)
		if !ok {
			http.Error(writer, "转换失败", http.StatusInternalServerError)
			return
		}

		// 创建一个新的客户端
		newClient := rs.NewClient()
		clientID := newClient.ID
		ch := newClient.Chan
		// 存储客户端到服务器
		rs.PutClient(newClient)
		slog.Info("新客户端连接", slog.String("clientID", clientID))

		// 关闭通道
		defer rs.RemoveClientByID(clientID)

		// SSE headers
		writer.Header().Set("Content-Type", "text/event-stream")
		writer.Header().Set("Cache-Control", "no-cache")
		writer.Header().Set("Connection", "keep-alive")

		req := c.Request
		for {
			select {
			case msg := <-ch:
				// 发送消息到客户端
				_, err := fmt.Fprintf(writer, "data: %s\n\n", msg)
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

	err := router.Run(addr)
	if err != nil {
		return
	}
}
