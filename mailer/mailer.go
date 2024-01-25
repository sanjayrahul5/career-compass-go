package mailer

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"runtime"
)

// SendMail sends a mail using SMTP server
func SendMail(mailTopic string, mail string, data any) {
	m := gomail.NewMessage()

	m.SetHeader("From", config.TransportEmail)
	m.SetHeader("To", mail)

	switch mailTopic {
	case config.MailOTP:
		m.SetHeader("Subject", "Career Compass - OTP Authentication")
		m.SetBody("text/plain", fmt.Sprintf("Your OTP for signup is %s", data.(string)))
	}

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		config.TransportEmail,
		config.TransportEmailPassword,
	)

	err := dialer.DialAndSend(m)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error sending mail -> %s", err.Error()))
		return
	}

	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), "Mail sent successfully!")
}
