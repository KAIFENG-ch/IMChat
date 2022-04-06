package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// Manager 建立连接用户结构体
var Manager = &ClientManage{
	Clients:        ClientMap{Clients: make(map[string]*Client)},
	GroupBroadcast: make(chan *GroupBroadcast),
	Broadcast:      make(chan *Broadcast),
	Register:       make(chan *Client),
	Reply:          make(chan *Client),
	Unregister:     make(chan *Client),
}

// Connect 服务器与用户连接与断开连接
func (m *ClientManage) Connect() {
	for true {
		select {
		// 连接
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients.Lock()
			Manager.Clients.Clients[conn.ID] = conn
			Manager.Clients.Unlock()
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			_ = conn.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
			// 断开连接
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			Manager.Clients.Lock()
			if _, ok := Manager.Clients.Clients[conn.ID]; ok {
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected"})
				_ = conn.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
				close(conn.Send)
				//close(Manager.Broadcast)
				delete(Manager.Clients.Clients, conn.ID)
			}
			Manager.Clients.Unlock()
		case message := <-Manager.Broadcast:
			sendId := message.Client.SendID
			flag := false
			//_ = json.Unmarshal(message, &MessageStruct)
			Manager.Clients.Lock()
			for id, conn := range Manager.Clients.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message.Message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients.Clients, conn.ID)
				}
				Manager.Clients.Unlock()
			}
			if flag {
				log.Println("对方在线应答")
				replyMsg := &ReplyMsg{
					Code:    30000,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = message.Client.Socket.WriteMessage(websocket.TextMessage, msg)
			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    30001,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = message.Client.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		case message := <-Manager.GroupBroadcast:
			log.Printf("群消息已发送：%s", message.Message.Content)
			replyMsg := &ReplyMsg{
				Code:    30000,
				Content: "群消息已发送",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = message.Send.Socket.WriteMessage(websocket.TextMessage, msg)
			for _, c := range message.Client {
				if c == nil {
					continue
				}
				Manager.Clients.Lock()
				conn, ok := Manager.Clients.Clients[c.ID]
				Manager.Clients.Unlock()
				if ok {
					message.Message.SendID = message.Send.ID
					conn.Send <- message.Message
				}
			}
		}
	}
}
