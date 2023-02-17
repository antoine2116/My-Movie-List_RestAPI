package common

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB(uri string) {
	// Create mongo client
	ctx := context.TODO()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	// Connect to the database
	if err != nil {
		log.Fatalf("Error connecting to the database : %s", err)
	}

	// Ping database
	if client.Ping(ctx, nil); err != nil {
		log.Fatalf("Unable to ping database : %s", err)
	}

	log.Println("Connected to DB")

	DB = client.Database("apous-films")
}

func GetDB() *mongo.Database {
	return DB
}
