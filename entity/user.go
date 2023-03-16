package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"passwordHash,omitempty"`
	Provider     string             `bson:"provider,omitempty"`
}
