package router

import (
	"github.com/gin-gonic/gin"
	"michatroom/conf"
	"michatroom/model"
	"net/http"
)

func ManRouter() {
	hub := model.NewHub()
	go hub.Run()
	r := gin.Default()
	// middleware
	//Recovery中间件会恢复(recovers)任何恐慌(panics)如果存在恐慌，中间件将会写入500。这个中间件还是很有必要的，因为当你的程序里有异常情况你没考虑到的时候，程序就退出了，服务停止了，所以有必要的
	r.Use(gin.Recovery(), gin.Logger())
	// load static
	r.LoadHTMLGlob("public/templete/**/*")
	r.Static("assets/", "public/assets/")
	r.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/login.html", gin.H{
			"title": "聊s",
		})
	})
	// group
	usersRouter(r)
	chatRouter(r, hub)
	wsRouter(r, hub)
	// 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"code": http.StatusNotFound,
			"msg":  "找不到页面",
		})
	})

	// start server
	r.Run(conf.HttpPort)
}
