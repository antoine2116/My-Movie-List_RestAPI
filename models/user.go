package models

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"passwordHash"`
}

type RegisterUser struct {
	Username             string `json:"username" binding:"required"`
	Email                string `json:"email" binding:"required"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(regUser *RegisterUser) error {
	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	user := User{}
	user.Username = regUser.Username
	user.Email = regUser.Email
	user.PasswordHash = utils.HashPassword(regUser.Password)

	// Insert user
	if _, err := coll.InsertOne(ctx, &user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("User already exists with the same email")
		}
		panic(err)
	}

	// Create unique index for email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: -1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := coll.Indexes().CreateOne(ctx, indexModel); err != nil {
		panic(err)
	}

	return nil
}

func FindUserByEmail(email string) (*User, error) {
	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	var user *User

	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		panic(err)
	}

	return user, nil
}
