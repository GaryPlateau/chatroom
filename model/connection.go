package model

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"michatroom/utils"
	"time"
)

type Data struct {
	UUId    string
	GroupId string
	Type    string
	ToUuid  string
}

type Connection struct {
	Conn    *websocket.Conn
	Message chan []byte
	Data    *Data
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(3600 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Message:
			//SetWriteDeadline设置基础网络上的写入截止时间
			//联系写入超时后，websocket状态已损坏
			//所有未来的写入都将返回错误。t的零值表示写入不会超时。
			//c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				//发送消息关闭websocket消息
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			utils.ErrorHandler("创建websocket写失败", err)
			w.Write(message)
			n := len(c.Message)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Message)
			}
			err = w.Close()
			utils.ErrorHandler("关闭websocket失败", err)
		case <-ticker.C:
			//c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			//c.Conn.WriteMessage(websocket.TextMessage, []byte(c.Data.UUId+":"+"超时为发送信息，断开连接..."))
			err := c.Conn.WriteMessage(websocket.PingMessage, nil)
			utils.ErrorHandler("ping websocket失败", err)
		}
	}
}

func (c *Connection) ReadPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(1024)
	//c.Conn.SetReadDeadline(time.Now().Add(360 * time.Second))
	//c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		msgType, message, err := c.Conn.ReadMessage()
		if err != nil || msgType == websocket.CloseMessage {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.ErrorHandler("websocket读取错误", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		if c.Data.Type == "single" {
			SingleMsgSend(h, c.Data.UUId, c.Data.ToUuid, message, msgType)
		} else if c.Data.Type == "group" {
			GroupMsgSend(h, c.Data.GroupId, message, msgType)
		} else if c.Data.Type == "system" {
			BroadcastMsgSend(h, message, msgType)
		} else {
			fmt.Println("读取websocket信息错误!")
		}
	}
}
