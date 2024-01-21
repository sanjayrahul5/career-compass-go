package handlers

import (
	"career-compass-go/config"
	"career-compass-go/mailer"
	"career-compass-go/models"
	"career-compass-go/pkg/logging"
	"career-compass-go/service"
	"career-compass-go/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"runtime"
)

// Signup is the handler for new user registration
func Signup(c *gin.Context) {
	var form models.User

	err := c.ShouldBind(&form)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing the request body -> %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := service.CheckExistingDocument(config.UserCollection, bson.D{{"email", form.Email}})
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

	err = mailer.SendMail(config.MailOTP, form.Email, otp)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error sending otp mail to the user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"otp": otp}})
}

func Login(c *gin.Context) {

}
