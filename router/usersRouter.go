package router

import (
	"github.com/gin-gonic/gin"
	"michatroom/common"
	"michatroom/controller"
	"michatroom/utils"
)

func usersRouter(r *gin.Engine) {
	usersRouter := r.Group("login").Use(common.StartSession("varifyCode"), common.AccessJsMiddleware())
	{
		usersRouter.GET("/", controller.UsersController{}.UsersLogin)
		usersRouter.GET("/resetPwd", controller.UsersController{}.UsersResetPwd)
		usersRouter.GET("/register", controller.UsersController{}.UsersRegister)
		usersRouter.GET("/captcha", func(ctx *gin.Context) {
			utils.SetCaptcha(ctx, 4)
		})

		usersRouter.POST("/checkLogin", controller.UsersController{}.UsersLoginCheck)
		usersRouter.POST("/registerUser", controller.UsersController{}.UserRegister)
		usersRouter.POST("/checkUserExist", controller.UsersController{}.UsersExistCheck)
		usersRouter.POST("/resetPassword", controller.UsersController{}.UserResetPwd)
	}
}
