package users

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
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"passwordHash,omitempty"`
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

func CreateUser(regUser *RegisterUser) (User, error) {
	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	user := User{
		ID:           primitive.NewObjectID(),
		Username:     regUser.Username,
		Email:        regUser.Email,
		PasswordHash: utils.HashPassword(regUser.Password),
	}

	// Insert user
	if _, err := coll.InsertOne(ctx, &user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return user, errors.New("User already exists with the same email")
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

	return user, nil
}

func FindUserByEmail(email string) (User, error) {
	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	var user User

	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, err
		}
		panic(err)
	}

	return user, nil
}
