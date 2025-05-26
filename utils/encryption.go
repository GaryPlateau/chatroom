package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const SALT = "tifa"

func md5Encode(org string) (md5Str string) {
	has := md5.New()
	has.Write([]byte(org))
	dst := has.Sum(nil)
	md5Str = hex.EncodeToString(dst)
	return
}

func SaltingPwd(pwd string) (saltPwd string) {
	temp := pwd + SALT
	saltPwd = md5Encode(temp)
	return
}

func VerifyPwd(enterPwd, oldPwd string) bool {
	if enterPwd == oldPwd {
		return true
	} else {
		return false
	}
}

//html.UnescapeString
//html.EscapeString
