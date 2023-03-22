package users

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func Test_ProvidersGoogleExchange(t *testing.T) {
	asserts := assert.New(t)

	tokenEp := mockTokenEndpoint()
	defer tokenEp.Close()

	p := NewGoogleProvider("client-id", "client-secret", "redirect-uri")
	p.config.Endpoint.TokenURL = tokenEp.URL

	ctx := context.Background()

	// Successful
	token, err := p.Exchange(ctx, "authorization_code")
	asserts.Nil(err)
	asserts.Equal("AccessToken", token.AccessToken)

	// Fail
	_, err = p.Exchange(ctx, "invalid_code")
	asserts.NotNil(err)
}

func Test_ProvidersGoogleGetUserEmail(t *testing.T) {
	asserts := assert.New(t)

	profileEp := mockGoogleProfileEndpoint()
	defer profileEp.Close()

	p := NewGoogleProvider("client-id", "client-secret", "redirect-uri")
	p.profileUrl = profileEp.URL

	ctx := context.Background()

	// Successful
	email, err := p.GetUserEmail(ctx, &oauth2.Token{
		AccessToken: "access_token",
		Expiry:      time.Now().Add(1 * time.Minute),
	})
	asserts.Nil(err)
	asserts.Equal("steve@gmail.com", email)

	// Fail (invalid token)
	_, err = p.GetUserEmail(ctx, &oauth2.Token{
		AccessToken: "invalid_token",
		Expiry:      time.Now().Add(1 * time.Minute),
	})
	asserts.NotNil(err)
}

func Test_ProvidersGitHubExchange(t *testing.T) {
	asserts := assert.New(t)

	tokenEp := mockTokenEndpoint()
	defer tokenEp.Close()

	p := NewGitHubProvider("client-id", "client-secret", "redirect-uri")
	p.config.Endpoint.TokenURL = tokenEp.URL

	ctx := context.Background()

	// Successful
	token, err := p.Exchange(ctx, "authorization_code")
	asserts.Nil(err)
	asserts.Equal("AccessToken", token.AccessToken)

	// Fail
	_, err = p.Exchange(ctx, "invalid_code")
	asserts.NotNil(err)
}

func Test_ProvidersGitHubGetUserEmail(t *testing.T) {
	asserts := assert.New(t)

	emailsEp := mockGitHubEmailsEndpoint()
	defer emailsEp.Close()

	p := NewGitHubProvider("client-id", "client-secret", "redirect-uri")
	p.emailsUrl = emailsEp.URL

	ctx := context.Background()

	// Successful
	email, err := p.GetUserEmail(ctx, &oauth2.Token{
		AccessToken: "access_token",
		Expiry:      time.Now().Add(1 * time.Minute),
	})
	asserts.Nil(err)
	asserts.Equal("steve@gmail.com", email)

	// Fail (invalid token)
	_, err = p.GetUserEmail(ctx, &oauth2.Token{
		AccessToken: "invalid_token",
		Expiry:      time.Now().Add(1 * time.Minute),
	})
	asserts.NotNil(err)
}

func mockTokenEndpoint() *httptest.Server {
	token := &oauth2.Token{
		AccessToken: "AccessToken",
		Expiry:      time.Now().Add(1 * time.Minute),
	}

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			code := r.FormValue("code")

			// Failure
			if code == "invalid_code" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Success
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(token)
		},
	))

	return server
}

func mockGoogleProfileEndpoint() *httptest.Server {
	user := googleUser{
		ID:    "test1",
		Email: "steve@gmail.com",
	}

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			// Failure
			if authorization == "Bearer invalid_token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Success
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		},
	))

	return server
}

func mockGitHubEmailsEndpoint() *httptest.Server {
	emails := []gitHubEmail{
		{Email: "steve@gmail.com", Primary: true},
		{Email: "steve2@gmail.com", Primary: false},
	}

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")

			// Failure
			if authorization == "Bearer invalid_token" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Success
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(emails)
		},
	))

	return server
}
