package service

import (
	"IMChat/dao"
	"IMChat/model"
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
	//uid := c.Query("uid")
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
		//_, message, err := c.Socket.ReadMessage()
		message := new(SendMsg)
		err := c.Socket.ReadJSON(&message)
		if err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
		}
		switch message.Type {
		case 3:
			var users []model.User
			var group model.Group
			model.DB.Model(&model.Group{}).Where("id = ?", message.Group).First(&group)
			err := model.DB.Model(&group).Association("Users").Find(&users)
			if err != nil {
				panic(err)
			}
			var clients []*Client
			for _, u := range users {
				Manager.Clients.Lock()
				clients = append(clients, Manager.Clients.Clients[strconv.Itoa(int(u.ID))])
				Manager.Clients.Unlock()
			}
			broadcast := &GroupBroadcast{
				Send:    c,
				Client:  clients,
				Message: message,
			}
			Manager.GroupBroadcast <- broadcast
		case 2:
			var group model.Group
			var users []model.User
			model.DB.Model(&model.Group{}).Where("id = ?", message.Group).
				First(&group)
			_ = model.DB.Model(&group).Where("user_id = ?", c.ID).
				Association("Users").Find(&users)
			if len(users) == 0 {
				_ = c.Socket.WriteMessage(websocket.TextMessage, []byte("您不是群成员"))
				continue
			}
			broadcast := &Broadcast{
				Client: c,
				Message: &SendMsg{
					Type:  2,
					Group: message.Group,
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
			var user model.User
			var friends []model.User
			model.DB.Model(model.User{}).Where("id = ?", c.ID).First(&user)
			_ = model.DB.Model(&user).Where("friend_id = ?", c.SendID).
				Association("Friends").Find(&friends)
			if len(friends) == 0 {
				_ = c.Socket.WriteMessage(websocket.TextMessage,
					[]byte("该用户不是你的好友"))
				continue
			}
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
					var sendUser model.User
					model.DB.Model(&model.User{}).Where("id = ?", message.SendID).First(&sendUser)
					replyMsg := &ReplyMsg{
						From:    message.SendID,
						Content: sendUser.Name + "请求添加你为好友",
						Code:    30000,
					}
					msg, _ := json.Marshal(replyMsg)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				} else if message.Content == "1" {
					res := dao.PullGroup(c.ID, message.Group)
					msg, _ := json.Marshal(&res)
					_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
					log.Printf("发送到客户端的拉群申请")
				}
			case 1:
				res := dao.MakeFriends(c.ID, c.SendID)
				msg, _ := json.Marshal(&res)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				log.Printf("发送到客户端的申请")
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
