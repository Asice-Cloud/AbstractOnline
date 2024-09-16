package controller

import (
	"Chat/controller/websocket_work"
	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title": "Home",
	})
}

func Ws(c *gin.Context) {
	websocket_work.ServerWs(websocket_work.Global_Hub, c)
}
