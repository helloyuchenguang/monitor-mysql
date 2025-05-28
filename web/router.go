package web

import (
	"github.com/gin-gonic/gin"
	"monitormysql/global"
)

// StartServer 启动Web服务器
func StartServer() {
	// gin web
	router := gin.Default()
	router.Static("/static", "./static")

	// 设置index.html为默认页面
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// SSE路由
	router.GET("/sse", func(c *gin.Context) {
		sseRule := global.GetRule[[]byte](RuleName)
		sseRule.AddSSEClient(c.Writer, c.Request)
	})

	err := router.Run(":28080")
	if err != nil {
		return
	}
}
