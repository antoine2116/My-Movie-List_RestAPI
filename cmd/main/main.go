package main

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/config"
	"apous-films-rest-api/middlewares"
	"apous-films-rest-api/users"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	conf := config.LoadConfiguration("../../")

	// Router
	r := common.NewRouter()

	// Tests routes
	test := r.Group("/test")
	{
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	// Authentication routes
	auth := r.Group("/auth")
	{
		users.AddRoutes(auth)
	}

	// API
	api := r.Group("/api")
	api.Use(middlewares.JwtAuthentication())
	{
	}

	r.Run(fmt.Sprintf(":%v", conf.Server.Port))

	common.InitDB(conf.Database.URI)
}
