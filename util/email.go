package util

import (
	"github.com/jordan-wright/email"
	"gopkg.in/gomail.v2"
	"net/smtp"
	"pmis-backend-go/global"
	"strconv"
)

// SendEmail 使用email库，测试通过。更简单点
func SendEmail(to string, subject string, body string) error {
	host := global.Config.EmailConfig.OutgoingMailServer
	port := strconv.Itoa(global.Config.EmailConfig.Port)
	account := global.Config.EmailConfig.Account
	password := global.Config.EmailConfig.Password

	e := email.NewEmail()
	e.From = global.Config.EmailConfig.Account
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)
	err := e.Send(
		host+":"+port,
		smtp.PlainAuth("", account, password, host),
	)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}
	return nil
}

// SendEmail2 使用gomail库，测试通过
func SendEmail2(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "19725912@qq.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.qq.com", 465, "19725912@qq.com", "ejusnukrlniabgdd")
	err := d.DialAndSend(m)
	if err != nil {
		global.SugaredLogger.Errorln(err)
		return err
	}
	return nil
}
