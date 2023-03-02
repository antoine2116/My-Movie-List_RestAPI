package users

import (
	"apous-films-rest-api/common"
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
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"passwordHash,omitempty"`
}

func CreateUser(user *User) error {
	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

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

func FindUserById(user *User, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	db := common.GetDB()
	coll := db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()

	filter := bson.D{{Key: "_id", Value: objectId}}

	if err := coll.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		panic(err)
	}

	return nil
}
