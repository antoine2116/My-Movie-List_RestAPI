package main

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/database"
	"apous-films-rest-api/oauth"
	"apous-films-rest-api/router"
	"apous-films-rest-api/users"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load config
	cfg := config.Load("../../")

	// Connect to the database
	co := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.NewClient(co)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Printf("Database connection established")

	db := database.New(client.Database(cfg.Database.Dev))

	r := buildRouting(db, cfg)

	r.Run(fmt.Sprintf(":%v", cfg.Server.Port))
}

func buildRouting(db *database.DB, cfg *config.Config) *gin.Engine {
	// Router
	r := router.New()

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
		users.RegisterHandlers(auth,
			users.NewService(users.NewRepository(db), 
				cfg.Server.Secret, 
				cfg.Server.TokenDuration, 
				oauth.NewGoogleProvider(cfg.Google), 
				oauth.NewGitHubProvider(cfg.GitHub)),
			cfg.Client.URI,
		)
	}

	// API
	api := r.Group("/api")
	api.Use(users.JwtAuthentication(cfg.Server.Secret))
	{
	}

	return r
}
