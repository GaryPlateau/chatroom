package model

import (
	"fmt"
	"michatroom/driver"
	"michatroom/utils"
	"time"
)

type UserBasic struct {
	ID            int    `db:"id" gorm:"primarykey"`
	UUID          string `json:"uuid"`
	Username      string `valid:"matches(^[a-zA-Z][\\w]+$),length(8|16),required" json:"username"`
	Password      string `valid:"stringlength(8|16),required" json:"password"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$),required" json:"phone"`
	Email         string `valid:"email,required" json:"email"`
	Identify      string `json:"identity" db:"identify"`
	ClientIp      string `json:"clentIp" db:"client_ip"`
	ClientPort    string `json:"clientPort" db:"client_port"`
	LoginTime     int    `json:"loginTime" db:"login_time"`
	HeartbeatTime int    `json:"heartbeatTime" db:"heartbeat_time"`
	LogoutTime    int    `json:"logoutTime" db:"logout_time"`
	DeviceInfo    string `json:"deviceInfo" db:"device_info"`
	Photo         string `json:"photo"`
	Status        string `json:"status"`
	Nickname      string `json:"nickname"`
	CreatedAt     int    `json:"createdAt" db:"created_at"`
	UpdatedAt     int    `json:"updatedAt" db:"updated_at"`
	DeletedAt     int    `json:"deletedAt" db:"deleted_at"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

// 创建数据库表格
// db.AutoMigrate(&model.UserBasic{})
func CreateUserTable() {
	driver.MysqlSingleInstance.Db.AutoMigrate(&UserBasic{})
}

func CreateUser(user *UserBasic) {
	driver.MysqlSingleInstance.Db.Create(&user)
}

func DeleteUser(user *UserBasic) {
	driver.MysqlSingleInstance.Db.Delete(&user)
}

func UpdateUserByLoginTime(user *UserBasic) {
	driver.MysqlSingleInstance.Db.Model(&user).Updates(UserBasic{LoginTime: user.LoginTime, ClientIp: user.ClientIp})
}

func UpdateUserByPassword(user *UserBasic) {
	driver.MysqlSingleInstance.Db.Model(&user).Updates(UserBasic{Password: user.Password})
}

func FindUserByUUID(uuid string) (user *UserBasic) {
	driver.MysqlSingleInstance.Db.Where("uuid = ?", uuid).First(&user)
	return
}

func FindUserByNameOrPhoneOrEmail(username, phone, email string) (user *UserBasic) {
	driver.MysqlSingleInstance.Db.Where("username = ?", username).Or("phone = ?", phone).Or("email = ?", email).First(&user)
	return
}

func FindUserByPhone(phone string) (user *UserBasic) {
	driver.MysqlSingleInstance.Db.Where("phone = ?", phone).First(&user)
	return
}

func FindUserByEmail(email string) (user *UserBasic) {
	driver.MysqlSingleInstance.Db.Where("email = ?", email).First(&user)
	return
}

func CheckUserByUserNameAndPwd(username string, password string) (user *UserBasic) {
	nowTime := fmt.Sprintf("%d", time.Now().UnixMilli())
	enterPwd := utils.SaltingPwd(password)
	driver.MysqlSingleInstance.Db.Where("username = ? AND password = ?", username, enterPwd).First(&user)
	driver.MysqlSingleInstance.Db.Model(&user).Where("id = ?", user.ID).Update("login_time", nowTime)
	return
}
