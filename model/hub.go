package model

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	Lock                 sync.Mutex
	Clients              map[string]*Connection
	Groups               map[string][]string
	Register, Unregister chan *Connection
	SingleMessage        chan *SingleMessageData
	GroupMessage         chan *GroupMessageData
	BroadcastMessage     chan *BroadcastMessageData
}

func NewHub() *Hub {
	return &Hub{
		Clients:          make(map[string]*Connection, 1024),
		Groups:           make(map[string][]string, 128),
		Register:         make(chan *Connection),
		Unregister:       make(chan *Connection),
		SingleMessage:    make(chan *SingleMessageData),
		GroupMessage:     make(chan *GroupMessageData),
		BroadcastMessage: make(chan *BroadcastMessageData),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case connection := <-h.Register:
			h.Lock.Lock()
			h.Clients[connection.Data.UUId] = connection
			h.Lock.Unlock()
		case connection := <-h.Unregister:
			h.Lock.Lock()
			if _, ok := h.Clients[connection.Data.UUId]; ok {
				delete(h.Clients, connection.Data.UUId)
				close(connection.Message)
				if len(h.Groups[connection.Data.GroupId]) == 0 {
					delete(h.Groups, connection.Data.GroupId)
				}
			}
			h.Lock.Unlock()
			//单对单
		case singleMsg := <-h.SingleMessage:
			if client, ok := h.Clients[singleMsg.ToId]; ok {
				sMsgJson, _ := json.Marshal(singleMsg)
				client.Message <- sMsgJson
			}
			//群聊
		case groupMsg := <-h.GroupMessage:
			if clients, ok := h.Groups[groupMsg.GroupId]; ok {
				for _, client := range clients {
					if _, ok := h.Clients[client]; ok {
						h.Clients[client].Message <- groupMsg.Message
					}
				}
			}
			//广播
		case broadcastMsg := <-h.BroadcastMessage:
			for cliendId, cConnect := range h.Clients {
				select {
				case cConnect.Message <- broadcastMsg.Message:
				default:
					close(cConnect.Message)
					delete(h.Clients, cliendId)
				}
			}
		}
	}
}
