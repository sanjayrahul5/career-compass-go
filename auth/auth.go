package auth

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"runtime"
	"time"
)

// GenerateToken generates a signed JWT token for the user session
func GenerateToken(userID, email string) (string, error) {
	secretKey := []byte(config.JWTSecret)

	//expiryTime := time.Now().Add(time.Minute * 30).Unix()
	expiryTime := time.Unix(1<<63-62135596801, 999999999).Unix() // Never expiring token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"email":  email,
			"exp":   expiryTime,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error creating signed JWT -> %s", err.Error()))
		return "", err
	}

	return tokenString, nil
}
