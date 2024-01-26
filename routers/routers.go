package routers

import (
	"career-compass-go/handlers"
	"career-compass-go/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetupRouter creates and configures a new Gin engine router
func SetupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.Logger())

	router.GET("/readyz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	router.POST("/signup", handlers.Signup)
	router.PUT("/signup/callback", handlers.SignupCallback)

	router.PUT("/reset-password", handlers.ResetPassword)

	router.POST("/signin", handlers.Login)

	return router
}
