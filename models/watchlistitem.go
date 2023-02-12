package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type WatchlistItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	WatchlistID primitive.ObjectID `bson:"watchlist_id"`
	MovieID     int                `bson:"movie_id"`
}
