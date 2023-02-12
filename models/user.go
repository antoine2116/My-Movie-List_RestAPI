package models

import (
	"apous-films-rest-api/common"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
	Email    string             `bson:"email"`
}

type RegisterUser struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// TODO new logic improvements
func CreateUser(user RegisterUser) error {
	db := common.GetDB()
	coll := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	// Check if username already exists
	filter := bson.D{{Key: "username", Value: user.Username}}
	count, err := coll.CountDocuments(ctx, filter)

	if err != nil {
		panic(err)
	}

	if count > 0 {
		return errors.New("User with same username already exists")
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		panic(err)
	}

	return nil
}
