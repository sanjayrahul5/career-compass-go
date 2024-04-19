package service

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

// GetAll gets the answer documents
func (an *Answer) GetAll(filters []bson.E) ([]Answer, error) {
	answers := make([]Answer, 0)

	cursor, err := config.AnswerCollection.Find(context.TODO(), bson.D(filters))
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error fetching answer documents -> %s", err.Error()))
		return nil, err
	}

	err = cursor.All(context.TODO(), &answers)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error decoding answer documents from cursor -> %s", err.Error()))
		return nil, err
	}

	return answers, nil
}
