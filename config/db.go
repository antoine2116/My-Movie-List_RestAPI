package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBClient(uri string) *mongo.Client {
	// Create mongo client
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatalf("Error creating database client : %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// Connect to the database
	if client.Connect(ctx); err != nil {
		log.Fatalf("Error connecting to the database : %s", err)
	}

	// Ping database
	if client.Ping(ctx, nil); err != nil {
		log.Fatalf("Unable to ping database : %s", err)
	}

	log.Println("Connected to DB")

	return client
}
