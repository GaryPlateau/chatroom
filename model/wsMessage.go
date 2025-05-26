package model

type MessageData struct {
	Message []byte
	Type    int
}

// 广播消息
type BroadcastMessageData struct {
	//
	MessageData
}

// 群组消息
type GroupMessageData struct {
	GroupId string
	MessageData
}

// 单独信息
type SingleMessageData struct {
	UId  string
	ToId string
	MessageData
}

//type Message struct {
//	Sender    string `json:"sneder,omitempty"`
//	Recipient string `json:"recipient,omitempty"`
//	Content   string `json:"content,omitempty"`
//}

func SingleMsgSend(h *Hub, sendId string, sendToId string, sendMsg []byte, sendMsgType int) {
	smd := &SingleMessageData{
		UId:  sendId,
		ToId: sendToId,
		MessageData: MessageData{
			Message: sendMsg,
			Type:    sendMsgType,
		},
	}
	h.SingleMessage <- smd
}

func GroupMsgSend(h *Hub, SendGroupId string, sendMsg []byte, sendMsgType int) {
	gmd := &GroupMessageData{
		GroupId: SendGroupId,
		MessageData: MessageData{
			Message: sendMsg,
			Type:    sendMsgType,
		},
	}
	h.GroupMessage <- gmd
}

func BroadcastMsgSend(h *Hub, sendMsg []byte, sendMsgType int) {
	bmd := &BroadcastMessageData{
		MessageData: MessageData{
			Message: sendMsg,
			Type:    sendMsgType,
		},
	}
	h.BroadcastMessage <- bmd
}
