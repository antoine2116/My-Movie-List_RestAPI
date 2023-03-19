package users

import (
	"apous-films-rest-api/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var errDB = errors.New("database error")

func Test_service_Register_and_Login(t *testing.T) {
	asserts := assert.New(t)

	s := NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	// Successful register
	token, err := s.Register(ctx, "example@mail.com", "pass")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// User Already exists
	_, err = s.Register(ctx, "example@mail.com", "pass")
	asserts.ErrorContains(err, "user already exists with the same email")

	// Successful login
	token, err = s.Login(ctx, "example@mail.com", "pass")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed login (unknown email)
	_, err = s.Login(ctx, "nobody@mail.com", "pass")
	asserts.ErrorContains(err, "invalid email or password")

	// Failed login (wrong password)
	_, err = s.Login(ctx, "example@mail.com", "wrong")
	asserts.ErrorContains(err, "invalid email or password")
}

func Test_service_Google_Login(t *testing.T) {
	asserts := assert.New(t)

	s := NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	// Successful authentication
	token, err := s.GoogleLogin(ctx, "valid_code")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed authentication
	_, err = s.GoogleLogin(ctx, "invalid_code")
	asserts.ErrorContains(err, "an error occured during Google authentication")
}

func Test_service_GitHub_Login(t *testing.T) {
	asserts := assert.New(t)

	s := NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	// Successful authentication
	token, err := s.GitHubLogin(ctx, "valid_code")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed authentication
	_, err = s.GitHubLogin(ctx, "invalid_code")
	asserts.ErrorContains(err, "an error occured during GitHub authentication")
}

func Test_service_performOAuth(t *testing.T) {
	asserts := assert.New(t)

	s := service{&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}}
	ctx := context.Background()

	// Successful OAuth authentication
	token, err := s.performOAuth(ctx, "example@mail.com", "some_provider")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed OAuth authentication
	_, err = s.performOAuth(ctx, "error@mail.com", "some_provider")
	asserts.NotNil(err)
}

func Test_service_generateJWT(t *testing.T) {
	asserts := assert.New(t)

	s := service{&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}}

	token := s.generateJWT("some_id", "example@mail.com")
	asserts.NotEmpty(token)
}

// Mocks
type mockRepository struct {
	items []*models.User
}

func (r mockRepository) GetById(ctx context.Context, id string) (*models.User, error) {
	for _, item := range r.items {
		if item.ID == id {
			return item, nil
		}
	}

	return nil, mongo.ErrNoDocuments
}

func (r mockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, item := range r.items {
		if item.Email == email {
			return item, nil
		}
	}

	return nil, mongo.ErrNoDocuments
}

func (r *mockRepository) Insert(ctx context.Context, user *models.User) (string, error) {
	if user.Email == "error@mail.com" {
		return "", errDB
	}

	user.ID = primitive.NewObjectID().Hex()
	r.items = append(r.items, user)
	return user.ID, nil
}

func (r mockRepository) Count(ctx context.Context) (int, error) {
	return len(r.items), nil
}

type mockProvider struct {
}

func (p mockProvider) GetUserEmail(ctx context.Context, code string) (string, error) {
	if code == "invalid_code" {
		return "", errors.New("")
	}

	return "example@mail.com", nil
}
