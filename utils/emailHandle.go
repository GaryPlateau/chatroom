package utils

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
)

/*
ReplyTo     []strings

	From        string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        []byte // Plaintext message (optional)
	HTML        []byte // Html message (optional)
	Sender      string // override From as SMTP envelope sender (optional)
	Headers     textproto.MIMEHeader
	Attachments []*Attachment
	ReadReceipt []string
*/
func SendMail(code, subject string, toEmail []string) {

	mailUserName := "272222798@qq.com" //邮箱账号
	mailPassword := "yvoptsdikrtnbhdc" //邮箱授权码
	addr := "smtp.qq.com:465"          //TLS地址
	host := "smtp.qq.com"              //邮件服务器地址

	e := email.NewEmail()
	e.From = "Get <272222798@qq.com>"
	e.To = toEmail
	e.Subject = subject //发送的主题
	e.HTML = []byte("验证码：<h3>" + code + "</h3>")
	err := e.SendWithTLS(addr, smtp.PlainAuth("", mailUserName, mailPassword, host),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	ErrorHandler("发送邮件错误", err)
}
