package connection

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type WebsocketClient struct {
	Id                string
	Conn              *websocket.Conn
	LastHeartbeatTime int64
	
	RoomId int64
	UserId int64
}

var (
	Clients sync.Map
)

type WsMessageType = int

const (
	JoinRoom      WsMessageType = 1
	LeaveRoom     WsMessageType = 2
	StartGame     WsMessageType = 3
	GameOver      WsMessageType = 4
	ChangeWord    WsMessageType = 5
	EliminateUser WsMessageType = 6
)

type WsMessage struct {
	Type    WsMessageType `json:"type"`
	Message interface{}   `json:"message"`
}

func SendMessageByRoomId(roomId int64, message WsMessage) {
	bytes, _ := json.Marshal(message)
	clients := GetClientsByRoomId(roomId)
	for _, client := range clients {
		client.Conn.WriteMessage(websocket.TextMessage, bytes)
	}
}

func SendMessageByUserId(userId int64, message WsMessage) {
	bytes, _ := json.Marshal(message)
	client := getClientByUserId(userId)
	client.Conn.WriteMessage(websocket.TextMessage, bytes)
}

func GetClientsByRoomId(roomId int64) (clients []WebsocketClient) {
	clients = make([]WebsocketClient, 0)
	Clients.Range(func(key, value interface{}) bool {
		client := value.(WebsocketClient)
		if client.RoomId == roomId {
			clients = append(clients, client)
		}
		return true
	})
	return
}

func getClientByUserId(userId int64) WebsocketClient {
	var client WebsocketClient
	Clients.Range(func(key, value interface{}) bool {
		client = value.(WebsocketClient)
		if client.UserId == userId {
			return false
		}
		return true
	})
	return client
}
