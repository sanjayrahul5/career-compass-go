package service

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"runtime"
	"time"
)

// Answer collection schema
type Answer struct {
	ID         primitive.ObjectID `json:"answerID" bson:"_id,omitempty"`
	QuestionID primitive.ObjectID `json:"questionID" bson:"question_id" binding:"required"`
	Content    string             `json:"content" bson:"content" binding:"required"`
	UserID     primitive.ObjectID `json:"userID" bson:"user_id"`
	UserName   string             `json:"userName" bson:"user_name"`
	CreatedAt  time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updated_at"`
}

// Add inserts a answer document
func (an *Answer) Add() error {
	res, err := config.AnswerCollection.InsertOne(context.TODO(), an)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error inserting new answer document -> %s", err.Error()))
		return err
	}

	an.ID = res.InsertedID.(primitive.ObjectID)
	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Created answerID -> %s", an.ID.Hex()))

	return nil
}