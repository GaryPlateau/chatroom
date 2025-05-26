package router

import (
	"github.com/gin-gonic/gin"
	"michatroom/controller"
	"michatroom/model"
)

func wsRouter(r *gin.Engine, hub *model.Hub) {
	wsRouter := r.Group("ws")
	{
		wsRouter.GET("/server", func(context *gin.Context) {
			controller.WsService(context, hub)
		})
		wsRouter.GET("/home", controller.Home)
	}
}
