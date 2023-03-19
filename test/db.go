package test

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/database"
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *database.DB

func DB(t *testing.T) *database.DB {
	if db != nil {
		return db
	}

	cfg, err := config.Load("../")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	co := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.NewClient(co)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		t.Error(err)
		t.FailNow()
	}

	db = database.New(client.Database(cfg.Database.Test))

	return db
}

func ResetCollections(t *testing.T, db *database.DB, tables ...string) {
	for _, table := range tables {
		err := db.Collection(table).Drop(context.Background())

		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
