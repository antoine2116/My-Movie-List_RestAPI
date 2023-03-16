package users

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/entity"
	"apous-films-rest-api/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Register(ctx context.Context, email string, password string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
}

type User struct {
	entity.User
}

type service struct {
	repo          Repository
	secret        string
	tokenDuration int
	googleCfg     config.GoogleConfig
	gitHubCfg     config.GitHubConfig
}

func NewService(repo Repository, secret string, tokenDuration int, googleCfg config.GoogleConfig, gitHubCfg config.GitHubConfig) Service {
	return service{repo, secret, tokenDuration, googleCfg, gitHubCfg}
}

func (s service) Register(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	_, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return "", errors.New("User already exists with the same email")
		}
	}

	// Insert user
	res, err := s.repo.Insert(ctx, entity.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: utils.HashPassword(password),
		Provider:     "local",
	})

	// Generate JWT token
	token := utils.GenerateJWT(res.InsertedID.(primitive.ObjectID).Hex())

	return token, nil
}

func (s service) Login(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("Invalid email or password")
		}
	}

	if err := utils.CompareHashAndPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("Invalid email or password")
	}

	// Generate JWT token
	token := utils.GenerateJWT(user.ID.Hex())

	return token, nil
}
