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

	// Routes that require token verification
	authRouter := router.Group("/")
	authRouter.Use(middlewares.VerifyToken())

	authRouter.GET("/role", handlers.GetAllRoles)
	authRouter.GET("/:id/role", handlers.GetRole)

	authRouter.GET("/skill", handlers.GetAllSkills)
	authRouter.GET("/:id/skill", handlers.GetSkill)

	authRouter.GET("/search", handlers.Search)

	authRouter.POST("/question", handlers.AddQuestion)
	authRouter.GET("/:id/question", handlers.GetQuestions)
	authRouter.PUT("/:id/question", handlers.UpdateQuestion)

	authRouter.POST("/answer", handlers.AddAnswer)

	// ML Routes
	router.POST("/predict", handlers.Predict)

	return router
}
