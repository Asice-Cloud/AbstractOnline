package websocket_work

import (
	"Abstract/config"
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	write_wait       = 10 * time.Second
	pong_wait        = 60 * time.Second
	ping_period      = (pong_wait * 9) / 10
	max_message_size = 512
)

type Client struct {
	room *Room
	conn *websocket.Conn
	send chan []byte
	Name string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (cli *Client) read() {
	defer func() {
		cli.room.unregister <- cli
		cli.conn.Close()
		leaveMessage := fmt.Sprintf("%s has left the chat room", cli.Name)
		cli.room.broadcast <- []byte(leaveMessage)
	}()
	cli.conn.SetReadLimit(max_message_size)
	cli.conn.SetReadDeadline(time.Now().Add(pong_wait))
	cli.conn.SetPongHandler(func(string) error {
		cli.conn.SetReadDeadline(time.Now().Add(pong_wait))
		return nil
	})
	for {
		_, message, err := cli.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				config.Lg.Error(fmt.Sprintf("error occurs%v", err))
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		if strings.HasPrefix(string(message), "@ai") {
			query := strings.TrimSpace(strings.TrimPrefix(string(message), "@ai"))
			ai_response := get_ai_response(query)
			cli.send <- []byte(fmt.Sprintf("Flandre Scarlet : %s", ai_response))
		} else {
			if string(message) == "#show" {
				var userList []string
				for key, value := range cli.room.clients {
					if value {
						if key.Name == cli.Name {
							userList = append(userList, key.Name+"(you)")
						} else {
							userList = append(userList, key.Name)
						}
					}
				}
				msg := fmt.Sprintf("#show:%s", strings.Join(userList, ","))
				cli.send <- []byte(msg)
			} else {
				cli.room.broadcast <- []byte(fmt.Sprintf("%s say: %s", cli.Name, message))
			}
		}
	}
}

func (cli *Client) write() {
	timer := time.NewTicker(ping_period)
	defer func() {
		timer.Stop()
		cli.conn.Close()
	}()
	for {
		select {
		case message, ok := <-cli.send:
			cli.conn.SetWriteDeadline(time.Now().Add(write_wait))
			if !ok {
				cli.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := cli.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			n := len(cli.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-cli.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-timer.C:
			cli.conn.SetWriteDeadline(time.Now().Add(write_wait))
			if err := cli.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServerWs(hub *Room, ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		config.Lg.Error(fmt.Sprintf("%v", err))
		return
	}

	baseName := getRandomName()
	uniqueName := getUniqueName(baseName)

	client := &Client{
		room: hub,
		conn: conn,
		send: make(chan []byte, 256),
		Name: uniqueName,
	}
	client.room.register <- client
	welcome := fmt.Sprintf("welcome %s join the chat room", client.Name)
	number := fmt.Sprintf("current number of people in the chat room is %d", len(client.room.clients)+1)
	go func() {
		client.room.broadcast <- []byte(welcome)
		client.room.broadcast <- []byte(number)
	}()

	go client.read()
	go client.write()
}
