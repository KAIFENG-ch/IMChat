package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

// Manager 建立连接用户结构体
var Manager = &ClientManage{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func creatId(uid, toUid string) string {
	return uid + "_" + toUid
}

// Connect 服务器与用户连接与断开连接
func (m *ClientManage) Connect() {
	for true {
		select {
		// 连接
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			_ = conn.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
			// 断开连接
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				jsonMessage, _ := json.Marshal(&Message{Content: "A socket has disconnected"})
				_ = conn.Socket.WriteMessage(websocket.TextMessage, jsonMessage)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case message := <-Manager.Broadcast:
			sendId := message.Client.SendID
			flag := false
			//_ = json.Unmarshal(message, &MessageStruct)
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message.Message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			//id := message.Client.ID
			if flag {
				log.Println("对方在线应答")
				replyMsg := &ReplyMsg{
					Code:    30000,
					Content: "对方在线应答",
				}
				msg, err := json.Marshal(replyMsg)
				_ = message.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//err = dao.InsertMsg("IMChat", sendId, string(message.Message), 1, int64(3*24*30*time.Hour))
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    30001,
					Content: "对方不在线应答",
				}
				msg, err := json.Marshal(replyMsg)
				_ = message.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//err = dao.InsertMsg("IMChat", sendId, string(message.Message), 0, int64(3*24*30*time.Hour))
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}
		}
	}
}
