package config

import (
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"log"
	"strconv"
)

var (
	ViperConfig *viper.Viper

	MongoClient     *mongo.Client
	MongoDBName     string
	MongoDBConn     *mongo.Database
	UserCollection  *mongo.Collection
	RoleCollection  *mongo.Collection
	SkillCollection *mongo.Collection

	Templates *template.Template

	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string

	JWTSecret string

	MLServerURL string
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

	// Initialize and parse the mailer template files
	Templates = template.Must(template.ParseGlob("mailer/templates/*.html"))

	MongoDBName = ViperConfig.GetString("DB_NAME")

	SMTPHost = ViperConfig.GetString("SMTP_HOST")
	SMTPPort, _ = strconv.Atoi(ViperConfig.GetString("SMTP_PORT"))
	SMTPEmail = ViperConfig.GetString("SMTP_EMAIL")
	SMTPPassword = ViperConfig.GetString("SMTP_PASSWORD")

	JWTSecret = ViperConfig.GetString("JWT_SECRET")

	MLServerURL = ViperConfig.GetString("ML_SERVER_URL")
}
