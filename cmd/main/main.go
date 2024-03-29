package main

import (
	"context"
	"fmt"
	"log"
	"mml-rest-api/internal/config"
	"mml-rest-api/internal/router"
	"mml-rest-api/internal/users"
	"mml-rest-api/pkg/database"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg, err := config.Load("../../config")

	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	// Connect to the database
	co := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.NewClient(co)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}()

	log.Printf("Database connection established")

	db := database.New(client.Database(cfg.Database.Db))

	// Router
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
				users.NewGoogleProvider(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL),
				users.NewGitHubProvider(cfg.GitHub.ClientID, cfg.GitHub.ClientSecret, cfg.GitHub.RedirectURL)),
			cfg.Client.URI,
		)
	}

	// API
	api := r.Group("/api")
	api.Use(users.JwtAuthentication(cfg.Server.Secret))
	{
		users.RegisterAuthenticatedHandlers(api)
	}

	return r
}
