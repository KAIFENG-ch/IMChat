package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Hub struct {
	C          map[*Connection]bool
	Data       chan []byte
	Register   chan *Connection
	Unregister chan *Connection
}

type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     int      `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"userList"`
}

type Connection struct {
	ws   *websocket.Conn
	sc   chan []byte
	data *Data
}

var h = Hub{
	C:          make(map[*Connection]bool),
	Register:   make(chan *Connection),
	Unregister: make(chan *Connection),
	Data:       make(chan []byte),
}

func HubWebsocket(c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:   512,
		WriteBufferSize:  512,
		HandshakeTimeout: 5 * time.Hour,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	w := &Connection{sc: make(chan []byte, 256), ws: conn, data: &Data{}}
	h.Register <- w

}

func (c *Connection) Write() {
	defer func() {
		h.Unregister <- c
		_ = c.ws.Close()
	}()
	for msg := range c.sc {
		_ = c.ws.WriteMessage(websocket.TextMessage, msg)
	}
}

var userList []string

func (c *Connection) Read() {
	for {
		c.ws.PongHandler()
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		_ = json.Unmarshal(msg, c.data)
		switch c.data.Type {
		case 1:
			c.data.User = c.data.Content
			c.data.From = c.data.User
			userList = append(userList, c.data.User)
			h.Data <- msg
			log.Printf("客户进入聊天室：%s", msg)
		case 0:
			_, msg, err := c.ws.ReadMessage()
			if err != nil {
				return
			}
			h.Data <- msg
			log.Printf("收到客户的消息：%s", msg)
		case -1:
			c.data.User = c.data.Content

		}
		h.Data <- msg
		log.Printf("收到客户的信息：%s", msg)
	}
}
