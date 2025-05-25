package main

import (
	"github.com/gin-gonic/gin"
	"monitormysql"
	"net/http"
)

type Client struct {
	Chan chan monitormysql.UpdateInfo
}

func main() {
	// monitor-mysql
	monitormysql.Run("./config/config-5900x.yml")
	// gin web
	router := gin.Default()
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {

		c.File("./static/index.html")
	})

	router.GET("/stream/:handlerType", func(c *gin.Context) {
		handlerType := c.Param("handlerType")
		if handler, ok := monitormysql.TableRegistry[handlerType]; ok {
			handler.AddSseClient(c.Writer, c.Request)
		} else {
			c.String(http.StatusBadRequest, "Handler not found")
		}
	})

	err := router.Run(":28080")
	if err != nil {
		return
	}
}
