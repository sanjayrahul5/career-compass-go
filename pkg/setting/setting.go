package setting

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime"
	"strings"
)

// Setup sets up the project dependency configurations
func Setup() {
	var err error

	// Connect to mongo cluster
	mongoUri := config.ViperConfig.GetString("MONGO_URI")

	config.MongoClient, err = ConnectToMongo(mongoUri)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Failed to establish connection to mongo -> %s", err.Error()))
		panic(err)
	}

	config.UserCollection = config.MongoClient.Database(config.MongoDBName).Collection(config.ViperConfig.GetString("USER_COLLECTION"))

	go CreateTTLIndexForUsers()
}

// ConnectToMongo establishes a client connection to the given mongoDB URI
func ConnectToMongo(uri string) (*mongo.Client, error) {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi))
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error connecting to mongo -> %s", err.Error()))
		return nil, err
	}

	return client, nil
}

// CloseMongoClient closes the mongo client connection
func CloseMongoClient(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error closing mongo client connection -> %s", err.Error()))
		panic(err)
	}
}

// Ping verifies the established mongo connection
func Ping(client *mongo.Client) error {
	err := client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error pinging mongo -> %s", err.Error()))
		return err
	}

	logging.Logger.Info("Connected to MongoDB...")
	return nil
}

// CreateTTLIndexForUsers creates TTL for otp expiration handling in users collection
func CreateTTLIndexForUsers() {
	indexName := "expire_at_1"
	expireAfterSeconds := int32(config.TTLIndexExpirySeconds)

	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "expire_at", Value: 1}},
		Options: options.Index().SetName(indexName).SetExpireAfterSeconds(expireAfterSeconds),
	}

	_, err := config.UserCollection.Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		if strings.Contains(err.Error(), "An equivalent index already exists with the same name but different options") {
			logging.Logger.Warning(utils.GetFrame(runtime.Caller(0)), "TTL index already exists for users collection... Skipping TTL index creation")
		} else {
			logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error creating user TTL index -> %s", err.Error()))
		}
	}
}
