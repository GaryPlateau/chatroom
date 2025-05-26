package model

import (
	"gorm.io/gorm"
	"michatroom/driver"
)

type UserFriends struct {
	gorm.Model
	UUID        string `json:"uuid"`
	FriendId    string `json:"friendId" db:"friend_id"`
	Relation    string `json:"relation"`
	Favorite    string `json:"favorite"`
	Nickname    string `json:"nickname"`
	Passthrough string `json:"passThrough" db:"pass_through"`
	Remark      string `json:"remark"`
}

func (table *UserFriends) TableName() string {
	return "user_friends"
}

func GetUserFriendsList(uuid string) (ufs []*UserFriends) {
	driver.MysqlSingleInstance.Db.Where("uuid = ?", uuid).Find(&ufs)
	return
}
