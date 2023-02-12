package app

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/config"
	"apous-films-rest-api/handlers"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *gin.Engine
	DB     *mongo.Client
}

func (a *App) Initialize(c config.Configuration) {
	// Router
	a.Router = gin.Default()

	// Recovers from any panics and write a 500 if there was one
	a.Router.Use(gin.Recovery())

	v1 := a.Router.Group("/api")

	user := v1.Group("/user")
	user.POST("/register", handlers.UserRegestration)

	// Tests routes
	test := v1.Group("/test")
	test.GET("/ping", handlers.Ping)

	// Database
	common.InitDB(c.Database.URI)
}

func (a *App) Run(port int) {
	a.Router.Run(fmt.Sprintf(":%d", port))
}
