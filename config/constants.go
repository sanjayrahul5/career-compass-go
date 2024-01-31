package config

import "time"

const (
	TTLIndexExpirySeconds = 10 * 60
	OTPExpiryTime         = time.Minute * 10

	OTPChars  = "1234567890"
	OPTLength = 6
	MailOTP   = "MailOTP"

	UserRole = "user"
)
