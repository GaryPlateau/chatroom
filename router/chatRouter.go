package router

import (
	"github.com/gin-gonic/gin"
	"michatroom/common"
	"michatroom/controller"
	"michatroom/model"
)

func chatRouter(r *gin.Engine, hub *model.Hub) {
	chatRouter := r.Group("index").Use(common.JWTAuthMiddleware())
	{
		//显示页面
		chatRouter.GET("", controller.ChatController{}.ChatIndex)
		// get获取信息
		chatRouter.GET("/get_user_friends_list", controller.ChatController{}.GetUserFriendsList)
		chatRouter.GET("/get_chat_list", controller.ChatController{}.GetChatListHTML)
		// post请求信息
		chatRouter.POST("/post_chat_object", func(ctx *gin.Context) {
			controller.ChatController{}.CreateChatObject(ctx, hub)
		})
		chatRouter.POST("/create_chat_list", controller.ChatController{}.CreateChatListHTML)
	}
}
