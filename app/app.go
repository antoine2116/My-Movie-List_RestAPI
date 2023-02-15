package app

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/config"
	"apous-films-rest-api/middlewares"
	"apous-films-rest-api/users"
	"fmt"
	"net/http"

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
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	// Authentication routes
	auth := a.Router.Group("/auth")
	{
		auth.POST("/register", users.UserRegestration)
		auth.POST("/login", users.UserLogin)
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
