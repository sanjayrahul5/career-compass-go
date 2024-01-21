package setting

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"runtime"
	"time"
)

// Setup sets up the project dependency configurations
func Setup() {
	var err error

	// Connect to mongo cluster
	mongoUri := config.ViperConfig.GetString("MONGO_URI")

	config.MongoClient, config.MongoCtx, config.MongoCancelFunc, err = ConnectToMongo(mongoUri)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Failed to establish connection to mongo -> %s", err.Error()))
		panic(err)
	}

	config.UserCollection = config.MongoClient.Database(config.MongoDBName).Collection(config.ViperConfig.GetString("USER_COLLECTION"))
}

// ConnectToMongo establishes a client connection to the given mongoDB URI
func ConnectToMongo(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi))
	if err != nil {
		defer cancel()
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error connecting to mongo -> %s", err.Error()))
		return nil, nil, nil, err
	}

	return client, ctx, cancel, nil
}

// CloseMongoClient closes the mongo client connection
func CloseMongoClient(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error closing mongo client connection -> %s", err.Error()))
			panic(err)
		}
	}()
}

// Ping verifies the established mongo connection
func Ping(client *mongo.Client, ctx context.Context) error {
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error pinging mongo -> %s", err.Error()))
		return err
	}

	logging.Logger.Info("Connected to Mongo...")
	return nil
}
