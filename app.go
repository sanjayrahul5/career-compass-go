package main

import (
	"career-compass-go/config"
	"career-compass-go/pkg/logging"
	"career-compass-go/pkg/setting"
	"career-compass-go/routers"
	"career-compass-go/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"net/http"
	"runtime"
)

func init() {
	logging.Setup()
	setting.Setup()
}

func main() {
	port := fmt.Sprintf(":%s", config.ViperConfig.GetString("PORT"))

	gin.SetMode(gin.ReleaseMode)

	defer setting.CloseMongoClient(config.MongoClient, config.MongoCtx, config.MongoCancelFunc)
	err := setting.Ping(config.MongoClient, config.MongoCtx)
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error connecting to mongo -> %s", err.Error()))
		return
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"API-Token", "authorization", "Access-Control-Allow-Origin", "content-type", "Origin", "X-Requested-With", "Accept"},
		AllowedMethods: []string{"GET", "PUT", "POST", "DELETE"},
		ExposedHeaders: []string{"API-Token-Expiry"},
		MaxAge:         5,
	})

	server := &http.Server{
		Addr:    port,
		Handler: c.Handler(routers.SetupRouter()),
	}

	logging.Logger.Info(fmt.Sprintf("Server listening at port %s", port))

	err = server.ListenAndServe()
	if err != nil {
		logging.Logger.Error(utils.GetFrame(runtime.Caller(0)), fmt.Sprintf("Error starting the server -> %s", err.Error()))
		return
	}
}
