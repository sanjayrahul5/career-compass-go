package utils

import (
	"career-compass-go/config"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"path/filepath"
	"runtime"
	"strings"
)

// GetFrame returns a formatted string representing the frame of the call
func GetFrame(function uintptr, file string, line int, _ bool) string {
	absPath, _ := filepath.Rel(strings.Split(file, "career-compass-go")[0]+"career-compass-go", file)

	arr := strings.Split(runtime.FuncForPC(function).Name(), ".")
	funcName := arr[len(arr)-1]
	if funcName == "0" {
		funcName = arr[len(arr)-1]
	}

	return fmt.Sprintf("[%s][%s][%d] ", absPath, funcName, line)
}

// GenerateOTP generates an otp of given length
func GenerateOTP(length int) (string, error) {
	buffer := make([]byte, length)

	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(config.OTPChars)
	for i := 0; i < length; i++ {
		buffer[i] = config.OTPChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}

// HashPassword creates an encrypted hash for the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPasswordHash verifies if the hash and the password string matches
func VerifyPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
