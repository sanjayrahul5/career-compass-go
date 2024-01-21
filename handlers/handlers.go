package handlers

import (
	"career-compass-go/auth"
	"career-compass-go/config"
	"career-compass-go/mailer"
	"career-compass-go/pkg/logging"
	"career-compass-go/service"
	"career-compass-go/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"runtime"
)

// Signup is the handler for new user registration
func Signup(c *gin.Context) {
	var user service.User

	err := c.ShouldBind(&user)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing the request body -> %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := user.CheckExistingUser()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error checking for existing user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if existingUser {
		logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), config.ExistingUserMsg)
		c.JSON(http.StatusConflict, gin.H{"error": config.ExistingUserMsg})
		return
	}

	otp, err := utils.GenerateOTP(config.OPTLength)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error generating otp -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = mailer.SendMail(config.MailOTP, user.Email, otp)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error sending otp mail to the user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"otp": otp}})
}

// Login is the handler for user login
func Login(c *gin.Context) {
	var user service.User

	err := c.ShouldBind(&user)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing the request body -> %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Get()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), config.UnauthorizedLoginMsg)
			c.JSON(http.StatusUnauthorized, gin.H{"error": config.UnauthorizedLoginMsg})
			return
		}

		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error checking for existing user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error generating JWT token -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token, "role": user.Role}})
}
