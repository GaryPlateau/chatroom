package model

import "michatroom/driver"

type ChatList struct {
	Id           int    `json:"id"`
	UUID         string `json:"uuid"`
	ToUuid       string `json:"friendId" db:"to_uuid"`
	LastMsgTime  int    `json:"lastTimeMsg" db:"last_time_msg""`
	LastMsg      string `json:"lastMsg" db:"last_msg"`
	NewMsgCount  int    `json:"newMsgCount" db:"new_msg_count"`
	OnlineStatus string `json:"onlineStatus" db:"online_status"`
	ReadedStatus string `json:"readedStatus" db:"readed_status"`
}

func (table *ChatList) TableName() string {
	return "chat_list"
}

func CreateChatList(uuid string, toUuid string, time int) {
	ml := &ChatList{
		UUID:        uuid,
		ToUuid:      toUuid,
		LastMsgTime: time,
	}
	driver.MysqlSingleInstance.Db.Create(&ml)
}

func FindToChatList(uuid, toUuid string) *ChatList {
	ml := &ChatList{}
	driver.MysqlSingleInstance.Db.Where("`uuid` = ? AND `to_uuid` = ?", uuid, toUuid).First(&ml)
	return ml
}

func FindAllChatList(uuid string) []*ChatList {
	mls := []*ChatList{}
	driver.MysqlSingleInstance.Db.Where("`uuid` = ?", uuid).Find(&mls)
	return mls
}
