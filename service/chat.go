package service

import (
	"IMChat/dao"
	"IMChat/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Client websocket用户
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan *SendMsg
}

type ClientMap struct {
	Clients map[string]*Client
	sync.Mutex
}

// ClientManage websocket管理
type ClientManage struct {
	Clients        ClientMap
	Broadcast      chan *Broadcast
	GroupBroadcast chan *GroupBroadcast
	Reply          chan *Client
	Register       chan *Client
	Unregister     chan *Client
}

type Broadcast struct {
	Client  *Client
	Message *SendMsg
}

type GroupBroadcast struct {
	GroupId int
	Send    *Client
	Client  []*Client
	Message *SendMsg
}

type SendMsg struct {
	SendID  string
	Type    int
	Content string
	Group   int
}

// ReplyMsg 回复信息
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

func WsHandler(c *gin.Context) {
	claims, _ := utils.ParseToken(c.GetHeader("Authorization"))
	uid := strconv.Itoa(int(claims.Id))
	if uid == "0" {
		return
	}
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
		log.Println(err)
	}
	client := &Client{
		ID:     uid,
		SendID: toUid,
		Socket: conn,
		Send:   make(chan *SendMsg),
	}
	Manager.Register <- client
	go client.Write()
	go client.Read()
}

func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for true {
		c.Socket.PongHandler()
		message := new(SendMsg)
		err := c.Socket.ReadJSON(&message)
		if err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
		}
		switch message.Type {
		case 3:
			if dao.IsMember(c.SendID+"group", strconv.Itoa(message.Group)) {
				users := dao.FindMembers(message.Group)
				var clients []*Client
				for _, u := range users {
					Manager.Clients.Lock()
					clients = append(clients, Manager.Clients.Clients[strconv.Itoa(int(u.ID))])
					Manager.Clients.Unlock()
				}
				broadcast := &GroupBroadcast{
					GroupId: message.Group,
					Send:    c,
					Client:  clients,
					Message: message,
				}
				Manager.GroupBroadcast <- broadcast
			} else {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("您不是群成员"))
			}
		case 2:
			users := dao.FindGroupUser(message.Group, c.ID)
			if len(users) == 0 {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("您不是群成员"))
				continue
			}
			broadcast := &Broadcast{
				Client: c,
				Message: &SendMsg{
					Type:    2,
					Content: message.Content,
					Group:   message.Group,
				},
			}
			log.Printf("收到用户的拉群申请")
			Manager.Broadcast <- broadcast
		case 1:
			broadcast := &Broadcast{
				Client: c,
				Message: &SendMsg{
					Type:    1,
					Content: message.Content,
				},
			}
			log.Printf("收到客户的申请")
			Manager.Broadcast <- broadcast
		case 0:
			if dao.IsMember(c.SendID, c.ID) {
				broadcast := &Broadcast{
					Client: c,
					Message: &SendMsg{
						Type:    message.Type,
						Content: message.Content,
					},
				}
				if string(message.Content) == "" {
					return
				}
				log.Printf("收到客户的信息:%s", message.Content)
				Manager.Broadcast <- broadcast
			}
		}
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
			switch message.Type {
			case 3:
				log.Printf("收到群聊的消息：%s", message.Content)
				replyMsg := &ReplyMsg{
					From:    message.SendID,
					Code:    200,
					Content: message.Content,
				}
				msg, _ := json.Marshal(&replyMsg)
				Manager.Clients.Lock()
				_ = Manager.Clients.Clients[c.ID].
					Socket.WriteMessage(websocket.TextMessage, msg)
				Manager.Clients.Unlock()
			case 2:
				if message.Content == "0" {
					sendUser := dao.FindUser(c.ID)
					group := dao.FindOneGroup(message.Group)
					replyMsg := &ReplyMsg{
						From:    c.ID,
						Content: sendUser.Name + " 请求你加入群聊" + group.Name,
					}
					msg, _ := json.Marshal(&replyMsg)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
					log.Printf("发送到客户端的拉群申请")
				} else if message.Content == "1" {
					res := dao.PullGroup(c.ID, message.Group)
					msg, _ := json.Marshal(&res)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
					log.Printf("发送到客户端的进群申请")
				}
			case 1:
				if message.Content == "0" {
					sendUser := dao.FindUser(c.ID)
					replyMsg := &ReplyMsg{
						From:    c.ID,
						Content: sendUser.Name + "请求添加你为好友",
						Code:    30000,
					}
					msg, _ := json.Marshal(replyMsg)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
					log.Printf("发送到客户端的加好友申请")
				} else if message.Content == "1" {
					res := dao.MakeFriends(c.ID, c.SendID)
					msg, _ := json.Marshal(&res)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
					log.Printf("发送到客户端的好友通过申请")
				}
			case 0:
				log.Printf("发送到客户端的消息:%s", message.Content)
				replyMsg := &ReplyMsg{
					From:    c.ID,
					Code:    200,
					Content: message.Content,
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
