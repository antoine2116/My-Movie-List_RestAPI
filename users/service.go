package users

import (
	"apous-films-rest-api/entity"
	"apous-films-rest-api/oauth"
	"apous-films-rest-api/utils"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Register(ctx context.Context, email string, password string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, error)
	GitHubLogin(ctx context.Context, code string) (string, error)
	GetById(ctx context.Context, id string) (User, error)
}

type User struct {
	entity.User
}

type service struct {
	repo           Repository
	secret         string
	tokenDuration  int
	googleProvider oauth.GoogleProvider
	gitHubProvider oauth.GitHubProvider
}

func NewService(repo Repository, secret string, tokenDuration int, googleProvider oauth.GoogleProvider, gitHubProvider oauth.GitHubProvider) Service {
	return service{repo, secret, tokenDuration, googleProvider, gitHubProvider}
}

func (s service) Register(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	_, err := s.repo.GetByEmail(ctx, email)

	if err == nil {
		return "", errors.New("user already exists with the same email")
	}

	// Insert user
	res, err := s.repo.Insert(ctx, entity.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		PasswordHash: utils.HashPassword(password),
		Provider:     "local",
	})

	if err != nil {
		return "", nil
	}

	// Generate JWT token
	token := s.generateJWT(res.InsertedID.(primitive.ObjectID).Hex(), email)

	return token, nil
}

func (s service) Login(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
	}

	if err := utils.CompareHashAndPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token := s.generateJWT(user.ID.Hex(), email)

	return token, nil
}

func (s service) GoogleLogin(ctx context.Context, code string) (string, error) {
	// Get google user email
	email, err := s.googleProvider.GetUserEmail(ctx, code)

	if err != nil {
		return "", errors.New("an error occured during Google authentication")
	}

	return s.performOAuth(ctx, email, "google")
}

func (s service) GitHubLogin(ctx context.Context, code string) (string, error) {
	// Get google user email
	email, err := s.gitHubProvider.GetUserEmail(ctx, code)

	if err != nil {
		return "", errors.New("an error occured during GitHub authentication")
	}

	return s.performOAuth(ctx, email, "github")
}

func (s service) GetById(ctx context.Context, id string) (User, error) {
	user, err := s.repo.GetById(ctx, id)

	if err != nil {
		return User{}, err
	}

	return User{*user}, nil
}

func (s service) performOAuth(ctx context.Context, email string, provider string) (string, error) {
	// Check if user already exsists
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil && err == mongo.ErrNoDocuments {
		// Create user
		res, err := s.repo.Insert(ctx, entity.User{
			ID:           primitive.NewObjectID(),
			Email:        email,
			PasswordHash: "",
			Provider:     provider,
		})

		if err != nil {
			return "", err
		}

		user.ID = res.InsertedID.(primitive.ObjectID)
	}

	token := s.generateJWT(user.ID.Hex(), email)

	return token, nil
}

func (s service) generateJWT(userId string, email string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * time.Duration(s.tokenDuration)).Unix(),
		"id":    userId,
		"email": email,
	}

	tokenString, err := token.SignedString([]byte(s.secret))

	if err != nil {
		panic(err)
	}

	return tokenString
}
