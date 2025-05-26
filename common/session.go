package common

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"michatroom/conf"
)

const KEY = "TifaLockhart"

func AuthSessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// extract session
		sessionID := ctx.Request.Header.Get("Sessionid")
		if sessionID == "" {
			ctx.JSON(401, gin.H{
				"code": 1,
				"msg":  "unauthorized",
			})
			ctx.Abort()
			return
		}
		// set simple var
		ctx.Set("sessionID", sessionID)
		// validate session：1.check cache 2.check database // TODO

		ctx.Next()
	}
}

// register and login will save session
func SetSession(ctx *gin.Context, key string, value string) {
	session := sessions.Default(ctx)
	session.Set(key, value)
	session.Save()
	// simple encrypt: sha256(usename + string(timestamp))
	//contact := username + strconv.FormatInt(time.Now().UnixMilli(), 10)
	//sum := sha256.Sum256([]byte(contact))
	//val := fmt.Sprintf("%x", sum)
	//sessionID := session.Get("sessionID").(string)
	//ctx.Writer.Header().Set("sessionID", sessionID)
	//return sessionID
}

func StartSession(key string) gin.HandlerFunc {
	store := cookie.NewStore([]byte(KEY))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   conf.MaxExpireTime,
		HttpOnly: true,
		//SameSite: 4
	})
	return sessions.Sessions(key, store)
}

func HasSession(ctx *gin.Context, key string) bool {
	session := sessions.Default(ctx)
	//if val := session.Get("sessionID"); val == nil {
	if val := session.Get(key); val == nil {
		return false
	}
	return true
}

func GetSession(ctx *gin.Context, key string) string {
	session := sessions.Default(ctx)
	//val := session.Get("sessionID")
	val := session.Get(key)
	if val.(string) == "" {
		return ""
	}
	return val.(string)
}

// logout will clear session
func ClearSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}
