package config

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var (
	ViperConfig *viper.Viper

	MongoClient     *mongo.Client

	MongoDBName    string
	UserCollection *mongo.Collection

	TransportEmail         string
	TransportEmailPassword string
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigName("app")
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	ViperConfig = viper.GetViper()

	MongoDBName = ViperConfig.GetString("DB_NAME")

	TransportEmail = ViperConfig.GetString("TRANSPORT_EMAIL")
	TransportEmailPassword = ViperConfig.GetString("TRANSPORT_EMAIL_PASSWORD")
}
