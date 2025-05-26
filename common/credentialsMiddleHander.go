package common

import (
	"github.com/gin-gonic/gin"
)

// gin中cookie跨域设置 跨域实现
func AccessJsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
		w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId")
		w.Header().Set("shower", "秀儿")
		/*
			origin := c.GetHeader("Origin")
			method := r.Method
				c.Header("Access-Control-Allow-Origin", origin) // 注意这一行，不能配置为通配符“*”号
				c.Header("Access-Control-Allow-Credentials", "true") // 注意这一行，必须设定为 true
				c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId") // 我们自定义的header字段都需要在这里声明
				c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
				c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control")
			// 放行所有OPTIONS方法
			if method == "OPTIONS" {
				//c.AbortWithStatus(http.StatusNoContent)
				c.AbortWithStatus(http.StatusOK)
			}
			// 处理请求
			c.Next()
		*/
		c.Next()
	}
}
