package test

import (
	"apous-films-rest-api/pkg/database"
	"context"
	"os"
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

	dbURI := os.Getenv("DB_URI")

	if dbURI == "" {
		t.FailNow()
	}

	co := options.Client().ApplyURI(dbURI)
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

	db = database.New(client.Database("test_db"))

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
