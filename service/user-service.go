package service

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"runtime"
	"time"
)

// User collection schema
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
	OTP      string             `bson:"otp,omitempty"`
	ExpireAt time.Time          `bson:"expire_at,omitempty"`
}

// RatingsData hold the assessment ratings data with ordered fields of a user
type RatingsData struct {
	RateWebDev          int `json:"rate_webDev"`
	RateInfraAutomation int `json:"rate_infraAutomation"`
	RateMLAlgo          int `json:"rate_mlAlgo"`
	RateUIUX            int `json:"rate_uiux"`
	RateBlockchain      int `json:"rate_blockchain"`
	RateAppDev          int `json:"rate_appDev"`
	RateCloud           int `json:"rate_cloud"`
	RateTesting         int `json:"rate_testing"`
	RateDataAnalytics   int `json:"rate_dataAnalytics"`
	RateEmbedded        int `json:"rate_embedded"`
	RateAI              int `json:"rate_ai"`
	RateCyber           int `json:"rate_cyber"`
	RateARVR            int `json:"rate_arvr"`
	RateCompArch        int `json:"rate_compArch"`
	RateNetwork         int `json:"rate_network"`
	RateProjectMan      int `json:"rate_projectMan"`
	RateGame            int `json:"rate_game"`
}

// Create inserts a new user document
func (us *User) Create() error {
	res, err := config.UserCollection.InsertOne(context.TODO(), us)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error inserting new user document -> %s", err.Error()))
		return err
	}

	us.ID = res.InsertedID.(primitive.ObjectID)
	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Created userID -> %s", us.ID.Hex()))

	return nil
}

// Get finds and returns the user document
func (us *User) Get(filters []bson.E) error {
	err := config.UserCollection.FindOne(context.TODO(), bson.D(filters)).Decode(us)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error finding user document -> %s", err.Error()))
		return err
	}

	return nil
}

// Update updates the user document based on the given update query
func (us *User) Update(filters []bson.E, update bson.D) error {
	_, err := config.UserCollection.UpdateOne(context.TODO(), bson.D(filters), update)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error updating user document -> %s", err.Error()))
		return err
	}

	return nil
}

// CheckExistingUser checks if a user document already exists
func (us *User) CheckExistingUser() (bool, error) {
	err := config.UserCollection.FindOne(context.TODO(), bson.D{{"email", us.Email}}).Decode(us)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error finding user document -> %s", err.Error()))
		return false, err
	}

	return true, nil
}
