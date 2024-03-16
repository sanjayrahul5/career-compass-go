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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"runtime"
	"time"
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

	// Check if the user already exists
	existingUser, err := user.CheckExistingUser()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error checking for existing user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if existingUser {
		logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("User already exists for email -> %s", user.Email))
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Encrypt password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error hashing password -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate OTP
	otp, err := utils.GenerateOTP(config.OPTLength)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error generating otp -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashedPassword
	user.OTP = otp
	user.Role = config.UserRole
	user.ExpireAt = time.Now().Add(config.OTPExpiryTime)

	// Create user with OTP and expiry time
	err = user.Create()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error creating new user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send OTP via email
	go mailer.SendMail(config.MailOTP, user.Email, otp)

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"userID": user.ID}})
}

// SignupCallback is the handler for otp verification and post user registration actions
func SignupCallback(c *gin.Context) {
	userID := c.Query("userID")
	otpEntered := c.Query("otp")

	// Convert userID hex to object
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing userID to object -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var user service.User
	user.ID = objectID

	// Check if the user document is expired
	filters := []bson.E{
		{"_id", objectID},
	}

	err = user.Get(filters)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("OTP expired for the userID -> %s", userID))
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP expired for the user"})
			return
		}

		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting user document -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if opt matches
	if user.OTP == otpEntered {
		// Remove otp and expiry field from the user's document
		updateFields := bson.D{
			{"$unset", bson.D{
				{"otp", ""},
				{"expire_at", ""},
			}},
		}

		err = user.Update(filters, updateFields)
		if err != nil {
			logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error removing otp and expiry fields from user document -> %s", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), "Entered otp does not match with the stored otp")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "The OTP entered is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "User registered successfully"})
}

// ResetPassword is the handler for resetting the user passwords
func ResetPassword(c *gin.Context) {
	var user service.User

	err := c.ShouldBind(&user)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing the request body -> %s", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPassword := user.Password

	// Check if the user exists
	existingUser, err := user.CheckExistingUser()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error checking for existing user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !existingUser {
		logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("No user registered with the email -> %s", user.Email))
		c.JSON(http.StatusNotFound, gin.H{"error": "User with this email does not exists"})
		return
	}

	// Encrypt new password
	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error hashing password -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update new password
	filters := []bson.E{
		{"email", user.Email},
	}

	updateFields := bson.D{
		{"$set", bson.D{
			{"password", newHashedPassword},
		}},
	}

	err = user.Update(filters, updateFields)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error updating the new resetted password -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "Password reset successful"})
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

	inputPassword := user.Password

	// Check if user exists
	filters := []bson.E{
		{"email", user.Email},
	}

	err = user.Get(filters)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("No user registered with the email -> %s", user.Email))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password entered is incorrect"})
			return
		}

		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error checking for existing user -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	// Validate the password
	isValid := utils.VerifyPasswordHash(inputPassword, user.Password)
	if !isValid {
		logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Failed password verification for user -> %s", user.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or password entered is incorrect"})
		return
	}

	// Generate auth token
	token, err := auth.GenerateToken(user.ID.Hex())
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error generating JWT token -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token, "role": user.Role, "username": user.Username, "email": user.Email}})
}

// GetAllRoles is the handler for fetching all the role details
func GetAllRoles(c *gin.Context) {
	var role service.Role

	// Get all role details
	response, err := role.GetAll()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting all the role details -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetRole is the handler for fetching role details
func GetRole(c *gin.Context) {
	roleID := c.Param("id")

	// Convert roleID hex to object
	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing roleID to object -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the role details
	filter := []bson.E{
		{"_id", objectID},
	}

	var role service.Role

	err = role.Get(filter)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting role details for role [%s] -> %s", roleID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": role})
}

// GetAllSkills is the handler for fetching all the skill details
func GetAllSkills(c *gin.Context) {
	var skill service.Skill

	// Get all skill details
	response, err := skill.GetAll()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting all the skill details -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}

// GetSkill is the handler for fetching skill details
func GetSkill(c *gin.Context) {
	skillID := c.Param("id")

	// Convert skillID hex to object
	objectID, err := primitive.ObjectIDFromHex(skillID)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error parsing skillID to object -> %s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the skill details
	filter := []bson.E{
		{"_id", objectID},
	}

	var skill service.Skill

	err = skill.Get(filter)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting skill details for skill [%s] -> %s", skillID, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": skill})
}
