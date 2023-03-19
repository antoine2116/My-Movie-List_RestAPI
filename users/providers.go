package users

import (
	"apous-films-rest-api/config"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	GetUserEmail(ctx context.Context, code string) (string, error)
}

type GoogleUser struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Verified   bool   `json:"verified_email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

type googleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider(cfg config.GoogleConfig) googleProvider {
	return googleProvider{
		&oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		},
	}
}

func (p googleProvider) GetUserEmail(ctx context.Context, code string) (string, error) {
	// Exchange will do the handshake to retrieve the initial access token.
	token, err := p.config.Exchange(ctx, code, oauth2.AccessTypeOffline)

	if err != nil {
		return "", err
	}
	user := GoogleUser{}

	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return "", err
	}

	return user.Email, err
}

type gitHubEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type gitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(cfg config.GitHubConfig) gitHubProvider {
	return gitHubProvider{
		&oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
	}
}

func (p gitHubProvider) GetUserEmail(ctx context.Context, code string) (string, error) {
	// Exchange will do the handshake to retrieve the initial access token.
	token, err := p.config.Exchange(ctx, code)

	if err != nil {
		return "", err
	}

	// Make a request to the GitHub API with to get the user's emails
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user/emails")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var emails []gitHubEmail

	err = json.NewDecoder(resp.Body).Decode(&emails)

	if err != nil {
		return "", err
	}

	// Return the user's primary email
	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", nil
}
