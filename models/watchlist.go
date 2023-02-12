package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Watchlist struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID string             `bson:"user_id"`
	Name   string             `bson:"name"`
}
