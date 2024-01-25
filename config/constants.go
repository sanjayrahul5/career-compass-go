package config

import "time"

const (
	TTLIndexExpirySeconds = 30 * 60
	OTPExpiryTime         = time.Minute * 30

	OTPChars  = "1234567890"
	OPTLength = 6
	MailOTP   = "MailOTP"

	ExistingUserMsg      = "User with this email already exists"
	UnauthorizedLoginMsg = "Email or password entered is incorrect"

	UserRole = "user"
)
