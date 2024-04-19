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

// Question collection schema
type Question struct {
	ID        primitive.ObjectID   `json:"questionID" bson:"_id,omitempty"`
	SkillID   primitive.ObjectID   `json:"skillID" bson:"skill_id" binding:"required"`
	Title     string               `json:"title" bson:"title" binding:"required"`
	Content   string               `json:"content" bson:"content" binding:"required"`
	Status    string               `json:"status" bson:"status"`
	UserID    primitive.ObjectID   `json:"userID" bson:"user_id"`
	UserName  string               `json:"userName" bson:"user_name"`
	Likes     int64                `json:"likes" bson:"likes"`
	LikedBy   []primitive.ObjectID `json:"likedBy" bson:"liked_by"`
	CreatedAt time.Time            `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt" bson:"updated_at"`
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
