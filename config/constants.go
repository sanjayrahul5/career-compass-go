package config

import "time"

const (
	TTLIndexExpirySeconds = 30 * 60
	OTPExpiryTime         = time.Minute * 30

	OTPChars  = "1234567890"
	OPTLength = 6
	MailOTP   = "MailOTP"

	UserRole = "user"
)
