package service

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"runtime"
)

// CheckExistingDocument checks if a document exists in the collection based on the given filter
func CheckExistingDocument(coll *mongo.Collection, filter bson.D) (bool, error) {
	var res any

	err := coll.FindOne(config.MongoCtx, filter).Decode(&res)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error finding document -> %s", err.Error()))
		return false, err
	}

	return true, nil
}
