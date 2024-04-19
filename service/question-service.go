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

// Question collection schema
type Question struct {
	ID        primitive.ObjectID   `json:"questionID" bson:"_id,omitempty"`
	SkillID   primitive.ObjectID   `json:"skillID" bson:"skill_id" binding:"required"`
	Title     string               `json:"title" bson:"title" binding:"required"`
	Content   string               `json:"content" bson:"content" binding:"required"`
	Status    string               `json:"status" bson:"status"`
	UserID    primitive.ObjectID   `json:"userID" bson:"user_id"`
	UserName  string               `json:"userName" bson:"user_name"`
	Upvote    int64                `json:"upvote" bson:"upvote"`
	UpvoteBy  []primitive.ObjectID `json:"upvoteBy" bson:"upvote_by"`
	CreatedAt time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated_at"`
	Answers   []Answer             `json:"answers,omitempty" bson:"-"`
}

// Add inserts a question document
func (qu *Question) Add() error {
	res, err := config.QuestionCollection.InsertOne(context.TODO(), qu)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error inserting new question document -> %s", err.Error()))
		return err
	}

	qu.ID = res.InsertedID.(primitive.ObjectID)
	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Created questionID -> %s", qu.ID.Hex()))

	return nil
}

// GetAll gets the question documents
func (qu *Question) GetAll(filters []bson.E) ([]Question, error) {
	questions := make([]Question, 0)

	cursor, err := config.QuestionCollection.Find(context.TODO(), bson.D(filters))
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error fetching question documents -> %s", err.Error()))
		return nil, err
	}

	err = cursor.All(context.TODO(), &questions)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error decoding question documents from cursor -> %s", err.Error()))
		return nil, err
	}

	return questions, nil
}

// Update updates fields of a specific question
func (qu *Question) Update(questionID primitive.ObjectID, update bson.D) error {
	_, err := config.QuestionCollection.UpdateByID(context.TODO(), questionID, update)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error updating question [%s] -> %s", questionID.Hex(), err.Error()))
		return err
	}

	return nil
}
