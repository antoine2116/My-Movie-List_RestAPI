package users

import (
	"apous-films-rest-api/models"
	"apous-films-rest-api/utils"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	Register(ctx context.Context, email string, password string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	GoogleLogin(ctx context.Context, code string) (string, error)
	GitHubLogin(ctx context.Context, code string) (string, error)
	GetById(ctx context.Context, id string) (*models.User, error)
}

type service struct {
	repo           Repository
	secret         string
	tokenDuration  int
	googleProvider OAuthProvider
	gitHubProvider OAuthProvider
}

func NewService(repo Repository, secret string, tokenDuration int, googleProvider OAuthProvider, gitHubProvider OAuthProvider) Service {
	return service{repo, secret, tokenDuration, googleProvider, gitHubProvider}
}

func (s service) Register(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	_, err := s.repo.GetByEmail(ctx, email)

	if err == nil {
		return "", errors.New("user already exists with the same email")
	}

	// Insert user
	id, err := s.repo.Insert(ctx, &models.User{
		Email:        email,
		PasswordHash: utils.HashPassword(password),
		Provider:     "local",
	})

	if err != nil {
		return "", nil
	}

	// Generate JWT token
	token := s.generateJWT(id, email)

	return token, nil
}

func (s service) Login(ctx context.Context, email string, password string) (string, error) {
	// Check if user exists
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := utils.CompareHashAndPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token := s.generateJWT(user.ID, email)

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

func (s service) GetById(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) performOAuth(ctx context.Context, email string, provider string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	// If user already exsists return the token
	if err == nil {
		return s.generateJWT(user.ID, user.Email), nil
	}

	// Otherwise create user
	id, err := s.repo.Insert(ctx, &models.User{
		Email:        email,
		PasswordHash: "",
		Provider:     provider,
	})

	if err != nil {
		return "", err
	}

	return s.generateJWT(id, email), nil
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
