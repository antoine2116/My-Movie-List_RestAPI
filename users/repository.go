package users

import (
	"apous-films-rest-api/database"
	"apous-films-rest-api/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id string) (*models.User, error)
	Insert(ctx context.Context, user *models.User) (string, error)
}

type repository struct {
	coll *mongo.Collection
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"passwordHash,omitempty"`
	Provider     string             `bson:"provider,omitempty"`
}

func NewRepository(db *database.DB) Repository {
	return repository{db.Collection("users")}
}

func (r repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	user := new(User)

	if err := r.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return toModel(user), nil
}

func (r repository) GetById(ctx context.Context, id string) (*models.User, error) {
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	user := new(User)

	if err := r.coll.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return toModel(user), nil
}

func (r repository) Insert(ctx context.Context, user *models.User) (string, error) {
	u := toUser(user)

	res, err := r.coll.InsertOne(ctx, &u)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func toModel(u *User) *models.User {
	return &models.User{
		ID:           u.ID.Hex(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Provider:     u.Provider,
	}
}

func toUser(u *models.User) *User {
	return &User{
		ID:           primitive.NewObjectID(),
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Provider:     u.Provider,
	}
}
