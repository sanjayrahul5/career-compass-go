package mailer

import (
	"bytes"
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

	m.SetHeader("From", config.SMTPEmail)
	m.SetHeader("To", mail)

	switch mailTopic {
	case config.MailOTP:
		var body bytes.Buffer

		err := config.Templates.ExecuteTemplate(&body, "otpTemplate.html", struct {
			OTP string
		}{
			OTP: data.(string),
		})

		if err != nil {
			logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error executing otp mail template -> %s", err.Error()))
			return
		}

		m.SetHeader("Subject", "Career Compass - OTP Authentication")
		m.SetBody("text/html", body.String())
	}

	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPEmail,
		config.SMTPPassword,
	)

	err := dialer.DialAndSend(m)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error sending mail -> %s", err.Error()))
		return
	}

	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), "Mail sent successfully!")
}
