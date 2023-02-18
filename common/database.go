package common

import (
	"apous-films-rest-api/config"
	"context"
	"flag"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func InitDB() {
	conf := config.LoadConfiguration("../")

	// Create mongo client
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Database.URI))

	// Connect to the database
	if err != nil {
		log.Fatalf("Error connecting to the database : %s", err)
	}

	// Ping database
	if client.Ping(ctx, nil); err != nil {
		log.Fatalf("Unable to ping database : %s", err)
	}

	log.Println("Connected to DB")

	if flag.Lookup("test.v") == nil {
		DB = client.Database(conf.Database.Dev)
	} else {
		DB = client.Database(conf.Database.Test)
	}
}

func InitTestDB() {
	InitDB()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	DB.Collection("users").DeleteMany(ctx, bson.M{})
}

func FreeTestDB() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	if err := DB.Drop(ctx); err != nil {
		panic(err)
	}
}

func GetDB() *mongo.Database {
	return DB
}
