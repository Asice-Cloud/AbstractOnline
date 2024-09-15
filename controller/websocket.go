package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// WebSocket
// @Summary	WebSocket
// @Tags WebSocket
// @Success	200	{string} reply
// @router /user/ws [get]
func WebSocketHandler(c *gin.Context) {
	// 获取WebSocket连接
	wsUpgrader := websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// 处理WebSocket消息
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			fmt.Printf("处理文本消息, %s\n", string(p))
			ws.WriteMessage(websocket.TextMessage, p)
			// c.Writer.Write(p)
		case websocket.BinaryMessage:
			fmt.Println("处理二进制消息")
		case websocket.CloseMessage:
			fmt.Println("关闭websocket连接")
			return
		case websocket.PingMessage:
			fmt.Println("处理ping消息")
			ws.WriteMessage(websocket.PongMessage, []byte("ping"))
		case websocket.PongMessage:
			fmt.Println("处理pong消息")
			ws.WriteMessage(websocket.PongMessage, []byte("pong"))
		default:
			fmt.Printf("未知消息类型: %d\n", messageType)
			return
		}
	}

}
