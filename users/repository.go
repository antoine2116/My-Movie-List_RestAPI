package users

import (
	"apous-films-rest-api/database"
	"apous-films-rest-api/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetById(ctx context.Context, id string) (*entity.User, error)
	Insert(ctx context.Context, user entity.User) (*mongo.InsertOneResult, error)
}

type repository struct {
	coll *mongo.Collection
}

func NewRepository(db *database.DB) Repository {
	return repository{db.Collection("users")}
}

func (r repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	user := entity.User{}

	if err := r.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r repository) GetById(ctx context.Context, id string) (*entity.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	user := entity.User{}

	if err := r.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r repository) Insert(ctx context.Context, user entity.User) (*mongo.InsertOneResult, error) {
	return r.coll.InsertOne(ctx, &user)
}
