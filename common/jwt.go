package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"michatroom/utils"
	"net/http"
	"time"
)

// 过期时间
const TokenExpireTime = 30 * time.Minute

// 设置jwt密钥secret
var MySecret = []byte("tifa&rinoa&yuna")

type MyClaims struct {
	Uuid string
	jwt.StandardClaims
}

// 生成token的函数
func GenerateToken(uuid string) (token string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(TokenExpireTime)
	claims := MyClaims{
		uuid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 设置token过期时间
			IssuedAt:  nowTime.Unix(),    //签发时间
			Issuer:    "SherryTysha",     // 设置jwt签发者
			Id:        "",                //按需求选这个，有些实现中, 会控制这个ID是不是在黑/白名单来判断是否还有效
			NotBefore: 0,                 //生效起始时间
			Subject:   "",                //主题
		},
	}
	// 生产token
	// 用指定的哈希方法创建签名对象
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 用上面的声明和签名对象签名字符串token
	// 1. 先对Header和PayLoad进行Base64URL转换
	// 2. Header和PayLoadBase64URL转换后的字符串用.拼接在一起
	// 3. 用secret对拼接在一起之后的字符串进行HASH加密
	// 4. 连在一起返回
	// 生成Token，指定签名算法和claims
	token, err = tokenClaims.SignedString(MySecret)
	utils.ErrorHandler("生成token签名错误", err)
	return
}

// ParseToken 验证token函数
func ParseToken(tokenStr string) (*MyClaims, error) {
	//对token密钥进行验证
	// 第三个参数: 提供一个回调函数用于提供要选择的秘钥, 回调函数里面的token参数,是已经解析但未验证的,可以根据token里面的值做一些逻辑, 如`kid`的判断
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	utils.ErrorHandler("验证token失败", err)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*MyClaims); ok && tokenClaims.Valid {
			return claims, err
		}
	}
	return nil, errors.New("token无法解析")
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetSameSite(http.SameSiteNoneMode)
		token, err := ctx.Cookie("token")
		if token == "" {
			ctx.HTML(http.StatusOK, "404.html", gin.H{
				"code": http.StatusForbidden,
				"msg":  "没有携带token",
			})
			ctx.Redirect(http.StatusMovedPermanently, "login")
			ctx.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token验证失败",
				"err":  err,
			})
			ctx.Redirect(http.StatusMovedPermanently, "login")
			ctx.Abort()
			return
		} else if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
			ctx.HTML(http.StatusOK, "404.html", gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "token已过期",
				"err":  err,
			})
			ctx.Redirect(http.StatusMovedPermanently, "login")
			ctx.Abort()
			return
		}
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		//authHeader := ctx.Request.Header.Get("Authorization")
		//if len(authHeader) == 0 {
		//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		//		Status: 401,
		//		Msg:    "token验证失败",
		//		Error:  "authorization不存在",
		//		Data:   "",
		//	})
		//	ctx.Abort()
		//	return
		//}
		//parts := strings.SplitN(authHeader, " ", 2)
		//if !(len(parts) == 2 && parts[0] == "Bearer") {
		//	ctx.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		//		Status: 401,
		//		Msg:    "token验证失败",
		//		Error:  "authHeader信息不对",
		//		Data:   "",
		//	})
		//	ctx.Abort()
		//	return
		//}
		//// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		//mc, err := ParseToken(parts[1])
		//if err != nil {
		//	ctx.JSON(http.StatusOK, Response{
		//		Status: 200,
		//		Msg:    "token验证成功",
		//		Error:  "",
		//		Data:   "",
		//	})
		//	ctx.Abort()
		//	return
		//}
		// 将当前请求的username信息保存到请求的上下文c上
		//ctx.Set("userId", claims.UserId)
		// 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
		//ctx.Next()
	}
}

func RenewToken(claims *MyClaims) (string, error) {
	if withinLimit(claims.ExpiresAt, 600) {
		return GenerateToken(claims.Uuid)
	}
	return "", errors.New("登录已过期")
}

func withinLimit(s int64, l int64) bool {
	e := time.Now().Unix()
	return e-s < l
}
