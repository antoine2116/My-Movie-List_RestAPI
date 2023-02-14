package app

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/config"
	"apous-films-rest-api/handlers"
	"apous-films-rest-api/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *gin.Engine
	DB     *mongo.Client
}

func (a *App) Initialize(c *config.Configuration) {
	// Router
	a.Router = gin.Default()

	// Recovers from any panics and write a 500 if there was one
	a.Router.Use(gin.Recovery())

	// Tests routes
	test := a.Router.Group("/test")
	{
		test.GET("/ping", handlers.Ping)
	}

	// Authentication routes
	auth := a.Router.Group("/auth")
	{
		auth.POST("/register", handlers.UserRegestration)
		auth.POST("/login", handlers.UserLogin)
	}

	// API
	api := a.Router.Group("/api")
	api.Use(middlewares.JwtAuthentication())
	{
	}

	// Database
	common.InitDB(c.Database.URI)
}

func (a *App) Run(port int) {
	a.Router.Run(fmt.Sprintf(":%d", port))
}
