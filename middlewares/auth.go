package middlewares

import (
	"career-compass-go/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// VerifyToken middleware verifies the validity and authenticity of a JWT token
func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing claims"})
			c.Abort()
			return
		}

		c.Set("userID", claims["userID"].(string))
		c.Set("email", claims["email"].(string))

		c.Next()
	}
}
