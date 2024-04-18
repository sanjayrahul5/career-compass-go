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
)

// Skill collection schema
type Skill struct {
	ID          primitive.ObjectID   `json:"skillID" bson:"_id,omitempty"`
	RoleIDs     []primitive.ObjectID `json:"-" bson:"role_ids"`
	Name        string               `json:"name" bson:"name"`
	Image       string               `json:"image" bson:"image"`
	Description string               `json:"description,omitempty" bson:"description"`
	Youtube     []string             `json:"youtube,omitempty" bson:"youtube"`
	Website     []string             `json:"website,omitempty" bson:"website"`
	Courses     []string             `json:"courses,omitempty" bson:"courses"`
	Roles       []Role               `json:"roles,omitempty"`
}

// Get gets the skill document based on the given filter
func (s *Skill) Get(filters []bson.E) error {
	err := config.SkillCollection.FindOne(context.TODO(), bson.D(filters)).Decode(s)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting skill document -> %s", err.Error()))
		return err
	}

	return nil
}

// GetAll gets all the skill documents
func (s *Skill) GetAll(filters []bson.E) ([]Skill, error) {
	skills := make([]Skill, 0)

	cursor, err := config.SkillCollection.Find(context.TODO(), bson.D(filters))
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting skill documents -> %s", err.Error()))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &skills)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error decoding skill documents from cursor -> %s", err.Error()))
		return nil, err
	}

	return skills, nil
}
