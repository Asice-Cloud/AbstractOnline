package user_module

import (
	"Chat/controller/websocket_work"
	"github.com/gin-gonic/gin"
)

// @Tags	home page
// @Success	200	{string} welcome
// @router	/user_module/home [get]
func Home(c *gin.Context) {
	c.HTML(200, "home.html", gin.H{
		"title": "Home",
	})
}

// @Tags	home page
// @Success	200	{string} welcome
// @router	/user_module/ws [get]
func Ws(c *gin.Context) {
	websocket_work.ServerWs(websocket_work.Global_Hub, c)
}
