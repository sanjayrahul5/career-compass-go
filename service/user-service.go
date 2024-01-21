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
)

// User collection schema
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
}

// Get finds and returns the user document
func (us *User) Get() error {
	err := config.UserCollection.FindOne(context.TODO(), bson.D{{"email", us.Email}}).Decode(us)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error finding user document -> %s", err.Error()))
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
