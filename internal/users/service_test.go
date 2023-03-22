package users

import (
	"apous-films-rest-api/internal/models"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

var errDB = errors.New("database error")

func Test_service_Register_and_Login(t *testing.T) {
	asserts := assert.New(t)

	s := NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	// Successful register
	token, err := s.Register(ctx, "steve@mail.com", "pass")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// User Already exists
	_, err = s.Register(ctx, "steve@mail.com", "pass")
	asserts.ErrorContains(err, "user already exists with the same email")

	// Successful login
	token, err = s.Login(ctx, "steve@mail.com", "pass")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed login (unknown email)
	_, err = s.Login(ctx, "notsteve@mail.com", "pass")
	asserts.ErrorContains(err, "invalid email or password")

	// Failed login (wrong password)
	_, err = s.Login(ctx, "steve@mail.com", "wrong")
	asserts.ErrorContains(err, "invalid email or password")
}

func Test_service_Google_Login(t *testing.T) {
	asserts := assert.New(t)

	s := NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	// Successful authentication (to test register)
	token, err := s.GoogleLogin(ctx, "valid_code")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// 2nd Successful authentication (to test login)
	token, err = s.GoogleLogin(ctx, "valid_code")
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

	// Successful authentication (to test register)
	token, err := s.GitHubLogin(ctx, "valid_code")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// 2nd Successful authentication (to test login)
	token, err = s.GitHubLogin(ctx, "valid_code")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed authentication
	_, err = s.GitHubLogin(ctx, "invalid_code")
	asserts.ErrorContains(err, "an error occured during GitHub authentication")
}

func Test_service_GetById(t *testing.T) {
	assert := assert.New(t)

	repo := &mockRepository{items: []*models.User{
		{ID: "test1", Email: "steve@gmail.com", PasswordHash: "hashed_password", Provider: "some_provider"},
	}}

	s := NewService(repo, "secret", 1000, mockProvider{}, mockProvider{})
	ctx := context.Background()

	user, err := s.GetById(ctx, "test1")
	assert.Nil(err)
	assert.Equal("steve@gmail.com", user.Email)

	_, err = s.GetById(ctx, "unknown_id")
	assert.NotNil(err)
}

func Test_service_performOAuth(t *testing.T) {
	asserts := assert.New(t)

	s := service{&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}}
	ctx := context.Background()

	// Successful OAuth authentication
	token, err := s.performOAuth(ctx, "steve@mail.com", "some_provider")
	asserts.Nil(err)
	asserts.NotEmpty(token)

	// Failed OAuth authentication
	_, err = s.performOAuth(ctx, "error@mail.com", "some_provider")
	asserts.NotNil(err)
}

func Test_service_generateJWT(t *testing.T) {
	asserts := assert.New(t)

	s := service{&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}}

	token := s.generateJWT("some_id", "steve@mail.com")
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

	return nil, errors.New("")
}

func (r mockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, item := range r.items {
		if item.Email == email {
			return item, nil
		}
	}

	return nil, errors.New("")
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

func (p mockProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	if code == "invalid_code" {
		return &oauth2.Token{}, errors.New("")
	}

	token := &oauth2.Token{
		AccessToken: "AccessToken",
		Expiry:      time.Now().Add(1 * time.Minute),
	}

	return token, nil
}

func (p mockProvider) GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	if token.AccessToken == "invalid_token" {
		return "", errors.New("")
	}

	return "steve@mail.com", nil
}
