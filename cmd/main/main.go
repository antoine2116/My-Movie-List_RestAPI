package main

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/config"
	"apous-films-rest-api/users"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	conf := config.LoadConfiguration("../../")

	// Router
	r := common.NewRouter()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		users.AddUserAuthentication(auth)
	}

	// API
	api := r.Group("/api")
	api.Use(users.JwtAuthentication())
	{
		users.AddUserProfile(api)
	}

	common.InitDB()

	r.Run(fmt.Sprintf(":%v", conf.Server.Port))
}
