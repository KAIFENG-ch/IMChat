package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// Client websocket用户
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// ClientManage websocket管理
type ClientManage struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// ReplyMsg 回复信息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// Message return结构体
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	conn, err := (&websocket.Upgrader{
		WriteBufferSize:  512,
		ReadBufferSize:   512,
		HandshakeTimeout: 5 * time.Hour,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &Client{
		ID:     creatId(uid, toUid),
		SendID: creatId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}
	Manager.Register <- client
	go client.Write()
	go client.Read()
	//defer func() {
	//	client.
	//}()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for true {
		c.Socket.PongHandler()
		_, message, _ := c.Socket.ReadMessage()
		broadcast := &Broadcast{
			Client:  c,
			Message: message,
		}
		log.Printf("收到客户的信息:%s", string(message))
		Manager.Broadcast <- broadcast
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for true {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到客户端的消息:%s", string(message))
			replyMsg := &ReplyMsg{
				Code:    200,
				Content: string(message),
			}
			msg, _ := json.Marshal(&replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
