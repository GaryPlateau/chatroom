package utils

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func captchaServer(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	var content bytes.Buffer
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Add("Access-Control-Allow-Headers", "Access-Token")
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		err := captcha.WriteImage(&content, id, width, height)
		ErrorHandler("验证码图片写入失败", err)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		err := captcha.WriteAudio(&content, id, lang)
		ErrorHandler("验证音频写入失败", err)
	default:
		return captcha.ErrNotFound
	}
	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

func SetCaptcha(ctx *gin.Context, length ...int) {
	l := captcha.DefaultLen
	w, h := 150, 39
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		l = length[1]
	}
	if len(length) == 3 {
		l = length[2]
	}
	if len(length) == 4 {
		l = length[3]
	}
	captchaId := captcha.NewLen(l)
	session := sessions.Default(ctx)
	session.Set("captcha", captchaId)
	err := session.Save()
	ErrorHandler("session保存失败", err)
	err = captchaServer(ctx.Writer, ctx.Request, captchaId, ".png", "zh", false, w, h)
}

func CaptchaVerify(c *gin.Context, code string) bool {
	session := sessions.Default(c)
	if captchaId := session.Get("captcha"); captchaId != nil {
		session.Delete("captcha")
		err := session.Save()
		ErrorHandler("删除captcha session失败", err)
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
