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

// Role collection schema
type Role struct {
	ID          primitive.ObjectID   `json:"roleID" bson:"_id,omitempty"`
	SkillIDs    []primitive.ObjectID `json:"skillIDs,omitempty" bson:"skill_ids"`
	Name        string               `json:"name" bson:"name"`
	Image       string               `json:"image" bson:"image"`
	Description string               `json:"description,omitempty" bson:"description"`
	Salary      string               `json:"salary,omitempty" bson:"salary"`
	Duties      []string             `json:"duties,omitempty" bson:"duties"`
	Companies   []Company            `json:"companies,omitempty" bson:"companies"`
	Skills      []Skill              `json:"skills,omitempty" bson:"-"`
}

type Company struct {
	Name  string `json:"name,omitempty" bson:"name"`
	Image string `json:"image,omitempty" bson:"image"`
	Link  string `json:"link,omitempty" bson:"link"`
}

// Create inserts a new role document
func (r *Role) Create() error {
	res, err := config.RoleCollection.InsertOne(context.TODO(), r)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error inserting role user document -> %s", err.Error()))
		return err
	}

	r.ID = res.InsertedID.(primitive.ObjectID)
	logging.Logger.Info(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Created roleID -> %s", r.ID.Hex()))

	return nil
}

// Get gets the role document based on the given filter
func (r *Role) Get(filters []bson.E) error {
	err := config.RoleCollection.FindOne(context.TODO(), bson.D(filters)).Decode(r)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting role document -> %s", err.Error()))
		return err
	}

	return nil
}

// GetAll gets all the role documents
func (r *Role) GetAll(filters []bson.E) ([]Role, error) {
	roles := make([]Role, 0)

	cursor, err := config.RoleCollection.Find(context.TODO(), bson.D(filters))
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error getting role documents -> %s", err.Error()))
		return nil, err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), &roles)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error decoding role documents from cursor -> %s", err.Error()))
		return nil, err
	}

	return roles, nil
}
