package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"michatroom/conf"
	"michatroom/driver"
	"michatroom/model"
	"michatroom/utils"
	"net/http"
)

type ChatController struct {
	BaseController
}

func (chatCntlr ChatController) ChatIndex(ctx *gin.Context) {
	var user model.UserBasic
	uuid := getClaims(ctx).Uuid
	temp := driver.RedisSingleInstance.HGetValue("users", uuid)
	err := json.Unmarshal([]byte(temp.(string)), &user)
	utils.ErrorHandler("聊天首页获取登录用户数据失败", err)

	ctx.HTML(http.StatusOK, "chat/index.html", gin.H{
		"title":    "聊s",
		"nickname": user.Nickname,
		"photo":    user.Photo,
		"host":     conf.HttpAddr + conf.HttpPort,
	})
}

// 获取好友列表
func (chatCntlr ChatController) GetUserFriendsList(ctx *gin.Context) {
	var user *model.UserBasic
	var userFriends []*model.UserFriends
	uuid := getClaims(ctx).Uuid

	if ufl := driver.RedisSingleInstance.HGetValue("userFriends", uuid); ufl.(string) == "" {
		userFriends = model.GetUserFriendsList(uuid)
		jUfl, _ := json.Marshal(userFriends)
		driver.RedisSingleInstance.HSetValue("userFriends", uuid, jUfl)
	} else {
		err := json.Unmarshal([]byte(ufl.(string)), &userFriends)
		utils.ErrorHandler("获取用户好友数据失败", err)
	}

	ufListHtml := ""
	for _, userfriend := range userFriends {
		if tuser := driver.RedisSingleInstance.HGetValue("users", userfriend.FriendId); tuser.(string) == "" {
			user = model.FindUserByUUID(userfriend.FriendId)
			juser, _ := json.Marshal(user)
			driver.RedisSingleInstance.HSetValue("users", uuid, juser)
		} else {
			err := json.Unmarshal([]byte(tuser.(string)), &user)
			utils.ErrorHandler("获取用户好友数据失败", err)
		}

		ufListHtml = ufListHtml + chatCntlr.createUserFriendsHTML(user)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"htmlData": ufListHtml,
	})
}

// 获取会话列表
func (chatCntlr ChatController) GetChatListHTML(ctx *gin.Context) {
	var beginLi, figure, userListBody, userListAction, result string
	var fUser *model.UserBasic
	var chatLists []*model.ChatList
	uuid := getClaims(ctx).Uuid

	if uuidCL := driver.RedisSingleInstance.HGetValue("chatList", uuid); uuidCL.(string) == "" {
		chatLists = model.FindAllChatList(uuid)
		jMls, _ := json.Marshal(chatLists)
		driver.RedisSingleInstance.HSetValue("chatList", uuid, jMls)
	} else {
		err := json.Unmarshal([]byte(uuidCL.(string)), &chatLists)
		utils.ErrorHandler("获取用户好友数据失败", err)
	}

	for _, chatlist := range chatLists {
		if tUser := driver.RedisSingleInstance.HGetValue("users", chatlist.ToUuid); tUser.(string) == "" {
			fUser = model.FindUserByUUID(chatlist.ToUuid)
			jfu, _ := json.Marshal(fUser)
			driver.RedisSingleInstance.HSetValue("users", uuid, jfu)
		} else {
			err := json.Unmarshal([]byte(tUser.(string)), &fUser)
			utils.ErrorHandler("获取用户好友数据失败", err)
		}

		beginLi = "<li id=\"" + fUser.UUID + "\" class=\"list-group-item\"><div>"
		if fUser.Photo == "" {
			figure = "<figure class=\"avatar\">\n<span class=\"avatar-title " + utils.RandBGStyle() + " rounded-circle\">" + string([]rune(fUser.Nickname)[0]) + "</span>\n</figure>"
		} else {
			figure = "<div>\n<figure class=\"avatar\">\n<img src=\"" + fUser.Photo + "\" class=\"rounded-circle\" alt=\"image\">\n</figure>\n</div>"
		}
		userListBody = "<div class=\"users-list-body\">\n<div>\n<h5>" + fUser.Nickname + "</h5>\n<p>" + chatlist.LastMsg + "</p>\n</div>"
		userListAction = "<div class=\"users-list-action\"><small class=\"text-muted\">" + utils.UnixToTime(chatlist.LastMsgTime) + "</small><div class=\"action-toggle\"><div class=\"dropdown\"><a data-toggle=\"dropdown\" href=\"#\"><i data-feather=\"more-horizontal\"></i></a><div class=\"dropdown-menu dropdown-menu-right\"><a href=\"#\" class=\"dropdown-item\">聊天</a><a href=\"#\" data-navigation-target=\"contact-information \"class=\"dropdown-item\">个人简介</a><a href=\"#\" class=\"dropdown-item\">添加收藏</a><div class=\"dropdown-divider\"></div><a href=\"#\" class=\"dropdown-item text-danger\">删除好友</a></div></div></div></div></div></li>"
		result = result + beginLi + figure + userListBody + userListAction
	}
	ctx.JSON(http.StatusOK, gin.H{
		"htmlData": result,
	})
	ctx.Next()
}

// 创建好友列表html
func (chatCntlr ChatController) createUserFriendsHTML(user *model.UserBasic) string {
	var figure string
	beginLi := "<li class=\"list-group-item\" id=\"" + user.UUID + "\">"
	if user.Photo == "" {
		figure = "<figure class=\"avatar\">\n<span class=\"avatar-title " + utils.RandBGStyle() + " rounded-circle\">" + string([]rune(user.Nickname)[0]) + "</span>\n</figure>"
	} else {
		figure = "<div>\n<figure class=\"avatar\">\n<img src=\"" + user.Photo + "\" class=\"rounded-circle\" alt=\"image\">\n</figure>\n</div>"
	}
	userListBody := "<div class=\"users-list-body\">\n<div>\n<h5>" + user.Nickname + "</h5>\n<p>" + user.Status + "</p>\n</div>"
	userListAction := "<div class=\"users-list-action\">\n<div class=\"action-toggle\">\n<div class=\"dropdown\">\n<a data-toggle=\"dropdown\" href=\"#\">\n<i data-feather=\"more-horizontal\"></i>\n</a>\n<div class=\"dropdown-menu dropdown-menu-right\">\n<a href=\"#\" class=\"dropdown-item\">消息</a>\n<a href=\"#\" data-navigation-target=\"contact-information\"\nclass=\"dropdown-item\">个人资料</a>\n<div class=\"dropdown-divider\"></div>\n<a href=\"#\" class=\"dropdown-item text-danger\">删除好友</a>\n</div>\n</div>\n</div>\n</div>"
	resultHTML := beginLi + figure + userListBody + userListAction + "</div></li>"
	return resultHTML
}

// 创建会话列表
func (chatCntlr ChatController) CreateChatListHTML(ctx *gin.Context) {
	uuid := getClaims(ctx).Uuid

	toUuid := ctx.PostForm("toUuid")
	if toUuid == "" || uuid == toUuid {
		ctx.JSON(http.StatusOK, gin.H{
			"htmlData": "",
		})
		ctx.Abort()
		return
	}
	mlHTML := chatCntlr.createChatItem(ctx, toUuid)
	ctx.JSON(http.StatusOK, gin.H{
		"htmlData": mlHTML,
	})
	ctx.Next()
}

// 创建与某个好友会话
func (chatCntlr ChatController) createChatItem(ctx *gin.Context, toUuid string) string {
	var figure, userListBody, userListAction string
	var fUser *model.UserBasic
	var chatlist *model.ChatList
	var chatlists []*model.ChatList
	uuid := getClaims(ctx).Uuid
	nowTime := int(utils.GetUnix())

	if toUuid == "" || toUuid == uuid {
		return ""
	}
	ml := model.FindToChatList(uuid, toUuid)
	if ml.Id == 0 {
		model.CreateChatList(uuid, toUuid, nowTime)
		chatlist.UUID = uuid
		chatlist.ToUuid = toUuid
		chatlist.LastMsgTime = nowTime

		if tUser := driver.RedisSingleInstance.HGetValue("users", toUuid); tUser.(string) == "" {
			fUser = model.FindUserByUUID(toUuid)
			jfu, _ := json.Marshal(fUser)
			driver.RedisSingleInstance.HSetValue("users", uuid, jfu)
		} else {
			err := json.Unmarshal([]byte(tUser.(string)), &fUser)
			utils.ErrorHandler("获取用户好友数据失败", err)
		}

		uChatList := driver.RedisSingleInstance.HGetValue("chatList", uuid)
		json.Unmarshal([]byte(uChatList.(string)), &chatlists)
		chatlists = append(chatlists, chatlist)
		jChatLists, _ := json.Marshal(chatlists)
		driver.RedisSingleInstance.HSetValue("chatList", uuid, jChatLists)

		beginLi := "<li class=\"list-group-item\"><div>"
		if fUser.Photo == "" {
			figure = "<figure class=\"avatar\">\n<span class=\"avatar-title " + utils.RandBGStyle() + " rounded-circle\">" + string([]rune(fUser.Nickname)[0]) + "</span>\n</figure>"
		} else {
			figure = "<div>\n<figure class=\"avatar\">\n<img src=\"" + fUser.Photo + "\" class=\"rounded-circle\" alt=\"image\">\n</figure>\n</div>"
		}
		userListBody = "<div id=\"" + toUuid + "\" class=\"users-list-body\">\n<div>\n<h5>" + fUser.Nickname + "</h5>\n<p>" + fUser.Status + "</p>\n</div>"
		userListAction = "<div class=\"users-list-action\"><small class=\"text-muted\">" + utils.UnixToTime(nowTime) + "</small><div class=\"action-toggle\"><div class=\"dropdown\"><a data-toggle=\"dropdown\" href=\"#\"><i data-feather=\"more-horizontal\"></i></a><div class=\"dropdown-menu dropdown-menu-right\"><a href=\"#\" class=\"dropdown-item\">聊天</a><a href=\"#\" data-navigation-target=\"contact-information \"class=\"dropdown-item\">个人简介</a><a href=\"#\" class=\"dropdown-item\">添加收藏</a><div class=\"dropdown-divider\"></div><a href=\"#\" class=\"dropdown-item text-danger\">删除好友</a></div></div></div></div></div></li>"
		resultHTML := beginLi + figure + userListBody + userListAction
		return resultHTML
	} else {
		return ""
		//userListBody = "<div class=\"users-list-body\">\n<div>\n<h5>" + user.Nickname + "</h5>\n<p>" + ml.LastMsg + "</p>\n</div>"
		//userListAction = "<div class=\"users-list-action\"><small class=\"text-muted\">" + ml.LastMsgTime + "</small><div class=\"action-toggle\"><div class=\"dropdown\"><a data-toggle=\"dropdown\" href=\"#\"><i data-feather=\"more-horizontal\"></i></a><div class=\"dropdown-menu dropdown-menu-right\"><a href=\"#\" class=\"dropdown-item\">聊天</a><a href=\"#\" data-navigation-target=\"contact-information \"class=\"dropdown-item\">个人简介</a><a href=\"#\" class=\"dropdown-item\">添加收藏</a><div class=\"dropdown-divider\"></div><a href=\"#\" class=\"dropdown-item text-danger\">删除好友</a></div></div></div></div></div></li>"
	}
}

// 创建聊天对象
func (chatCntlr ChatController) CreateChatObject(c *gin.Context, hub *model.Hub) {
	toUuid := c.PostForm("tid")
	if toUuid == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 400,
			"msg":    "建立链接失败!",
		})
		return
	}
	uuid := getClaims(c).Uuid
	conn, _ := hub.Clients[uuid]
	_, isOnline := hub.Clients[toUuid]
	if !isOnline {
		c.JSON(http.StatusOK, gin.H{
			"status": 400,
			"msg":    "当前用户不在线!",
		})
		return
	}
	conn.Data.Type = "single"
	conn.Data.ToUuid = toUuid
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "建立链接成功!",
	})
}
