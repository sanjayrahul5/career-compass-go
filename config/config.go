package config

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ViperConfig     *viper.Viper

	MongoClient     *mongo.Client
	MongoCtx        context.Context
	MongoCancelFunc context.CancelFunc
)

func init() {
}
