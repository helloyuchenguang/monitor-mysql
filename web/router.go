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

	router.GET("/", func(c *gin.Context) {

		c.File("./static/index.html")
	})

	router.GET("/sse", func(c *gin.Context) {
		//handlerType := c.Param("handlerType")
		sseRule := global.GetRule[[]byte]("SSERule")
		sseRule.AddSSEClient(c.Writer, c.Request)
	})

	err := router.Run(":28080")
	if err != nil {
		return
	}
}
