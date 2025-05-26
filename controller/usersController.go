package controller

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
	"michatroom/common"
	"michatroom/conf"
	"michatroom/driver"
	"michatroom/model"
	"michatroom/utils"
	"net/http"
)

type UsersController struct {
	BaseController
}

func (usersCntlr UsersController) UsersLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login/login.html", gin.H{
		"title": "聊s",
	})
}

func (usersCntlr UsersController) UsersResetPwd(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login/reset_password.html", gin.H{
		"title": "修改密码",
	})
}

func (usersCntlr UsersController) UsersRegister(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login/register.html", gin.H{
		"title": "注册用户",
	})
}

func (usersCntlr UsersController) UsersLoginCheck(ctx *gin.Context) {
	var user *model.UserBasic

	username := ctx.PostForm("jsonData[username]")
	password := ctx.PostForm("jsonData[password]")
	code := ctx.PostForm("jsonData[code]")

	if !utils.CaptchaVerify(ctx, code) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "输入的验证码有误,请重新输入!",
		})
		return
	}
	user = model.CheckUserByUserNameAndPwd(username, password)
	//driver.MysqlSingleInstance.Db.Model(&model.UserBasic{}).Where("username = ? AND password = ?", username, password).First(&user).Count(&count)
	if user.ID != 0 {
		token, err := common.GenerateToken(user.UUID)
		utils.ErrorHandler("获取token失败", err)
		user.ClientIp = ctx.ClientIP()
		user.ClientPort = ctx.RemoteIP()
		user.DeviceInfo = ctx.ContentType()
		user.LoginTime = int(utils.GetUnix())
		model.UpdateUserByLoginTime(user)

		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie("token", token, 3600, "/", conf.HttpAddr, false, true)

		jUser, _ := json.Marshal(user)
		driver.RedisSingleInstance.HSetValue("users", user.UUID, jUser)

		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    username + "登录成功",
			"url":    "index",
		})

	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "输入的账号或密码错误",
		})
		ctx.Abort()
		return
	}
	ctx.Next()
}

func (usersCntlr UsersController) UsersExistCheck(ctx *gin.Context) {
	originalAccount := ctx.PostForm("original_account")
	verifyCode := ctx.PostForm("verify_code")
	if !utils.CaptchaVerify(ctx, verifyCode) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusUnauthorized,
			"msg":    "输入的验证码有误,请重新输入!",
		})
		ctx.Abort()
		return
	}
	if originalAccount == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "请输入用户名或电话或邮箱",
		})
		ctx.Abort()
		return
	}
	user := model.FindUserByNameOrPhoneOrEmail(originalAccount, originalAccount, originalAccount)
	if user.ID != 0 {
		emailVerifyCode := utils.CreateRandString(6, 0)
		utils.SendMail(emailVerifyCode, "聊s修改密码提示", []string{user.Email})
		//ctx.SetSameSite(http.SameSiteNoneMode)
		session := sessions.Default(ctx)
		session.Set("emailVerifyCode", emailVerifyCode)
		session.Save()
		//ctx.SetCookie("emailVerifyCode", emailVerifyCode, 300, "/", conf.HttpAddr, false, false)
		ctx.JSON(http.StatusOK, gin.H{
			"status":   http.StatusOK,
			"htmlData": "<div class=\"form-group\">\n\t\t<input id=\"original_pwd\" type=\"text\" class=\"form-control\" placeholder=\"原密码\" required>\n\t</div><div class=\"form-group\">\n\t\t<input id=\"new_pwd\" type=\"text\" class=\"form-control\" placeholder=\"输入新密码\" required>\n\t</div><div class=\"form-group\">\n\t\t<input id=\"renew_pwd\" type=\"text\" class=\"form-control\" placeholder=\"再次输入新密码\" required>\n\t</div><div class=\"form-group\">\n\t\t<input id=\"email_verify_code\" type=\"text\" maxlength=\"6\" class=\"form-control\" placeholder=\"邮箱验证码\" required>\n\t</div>",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "请输入用户名或电话或邮箱",
		})
	}
	ctx.Next()
}

func (usersCntlr UsersController) UserResetPwd(ctx *gin.Context) {
	username := ctx.PostForm("username")
	originalPwd := ctx.PostForm("original_pwd")
	newPwd := ctx.PostForm("new_pwd")
	renewPwd := ctx.PostForm("renew_pwd")
	verifyCode := ctx.PostForm("email_verify_code")

	if username == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "更改密码失败，请稍后重试!",
		})
		ctx.Abort()
		return
	}

	user := model.FindUserByNameOrPhoneOrEmail(username, username, username)
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "用户不存在，修改密码失败!",
		})
		ctx.Abort()
		return
	}

	oPwd := utils.SaltingPwd(originalPwd)
	if oPwd != user.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "请输入密码错误",
		})
		ctx.Abort()
		return
	}

	if !utils.VerifyPwd(newPwd, renewPwd) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "两次输入的新密码不同",
		})
		ctx.Abort()
		return
	}

	session := sessions.Default(ctx)
	emailVerifyCode := session.Get("emailVerifyCode")
	//emailVerifyCode, err := ctx.Cookie("emailVerifyCode")
	//if err != nil {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"status": http.StatusForbidden,
	//		"msg":    "获取验证码失败!",
	//	})
	//	ctx.Abort()
	//	return
	//}
	if emailVerifyCode != verifyCode {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "输入的验证码有误,请重新输入!",
		})
		ctx.Abort()
		return
	}
	user.Password = utils.SaltingPwd(renewPwd)
	model.UpdateUserByPassword(user)
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "更改密码成功",
	})
	ctx.Next()
}

// 用户注册
func (usersCntlr UsersController) UserRegister(ctx *gin.Context) {
	valid := validator.New()
	var user *model.UserBasic
	username := ctx.PostForm("username")
	phone := ctx.PostForm("phone")
	email := ctx.PostForm("email")
	verifyCode := ctx.PostForm("verify_code")

	if !utils.CaptchaVerify(ctx, verifyCode) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "输入的验证码有误,请重新输入!",
		})
		ctx.Abort()
		return
	}

	if username == "" || phone == "" || email == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "请输入用户名或电话或邮箱为空",
		})
		ctx.Abort()
		return
	}

	user = model.FindUserByNameOrPhoneOrEmail(username, phone, email)

	if user.ID != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusForbidden,
			"msg":    "该用户名已经注册",
		})
		ctx.Abort()
		return
	}

	newPwd := utils.CreateRandString(8, 3)
	utils.SendMail(newPwd, "聊s注册账户密码", []string{email})
	newPwd = utils.SaltingPwd(newPwd)
	nowTime := int(utils.GetUnix())

	reqIP := ctx.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}

	newUser := &model.UserBasic{
		UUID:      uuid.NewV4().String(),
		Username:  username,
		Phone:     phone,
		Email:     email,
		Password:  newPwd,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
		ClientIp:  reqIP,
	}

	if err := valid.Struct(newUser); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusBadRequest,
			"msg":    err.Error(),
		})
		ctx.Abort()
		return
	}

	model.CreateUser(newUser)
	jnu, _ := json.Marshal(newUser)
	driver.RedisSingleInstance.HSetValue("users", newUser.UUID, jnu)

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "注册成功",
	})

	ctx.Next()
}

//data, err := c.GetRawData()
//common.ErrorHandler("打印出body multipart/form-data", err)
//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
//err = json.NewDecoder(c.Request.Body).Decode(&userVerifyInfo)
//jsonData[username]=haha&jsonData[password]=haha&jsonData[code]=1234
//jsonData, err := url.QueryUnescape(string(data))
//common.ErrorHandler("解析userVerifyInfo失败", err)
//arr := strings.Split(jsonData, "&")
//fmt.Println(arr)
