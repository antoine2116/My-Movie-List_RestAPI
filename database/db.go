package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	db *mongo.Database
}

func New(db *mongo.Database) *DB {
	return &DB{db}
}

func (db *DB) DB() *mongo.Database {
	return db.db
}

func (db *DB) Collection(name string) *mongo.Collection {
	return db.db.Collection(name)
}
