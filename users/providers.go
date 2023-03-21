package users

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type OAuthProvider interface {
	GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error)
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
}

type googleUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type googleProvider struct {
	config     *oauth2.Config
	profileUrl string
}

func NewGoogleProvider(clientID, clientSecret, redirectURL string) googleProvider {
	return googleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		},
		profileUrl: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

func (p googleProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	// Exchange will do the handshake to retrieve the initial access token.
	return p.config.Exchange(ctx, code, oauth2.AccessTypeOffline)
}

func (p googleProvider) GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	user := googleUser{}

	client := p.config.Client(ctx, token)
	resp, err := client.Get(p.profileUrl)

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
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

type gitHubProvider struct {
	config    *oauth2.Config
	emailsUrl string
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) gitHubProvider {
	return gitHubProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
		emailsUrl: "https://api.github.com/user/emails",
	}
}

func (p gitHubProvider) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	// Exchange will do the handshake to retrieve the initial access token.
	return p.config.Exchange(ctx, code, oauth2.AccessTypeOffline)
}

func (p gitHubProvider) GetUserEmail(ctx context.Context, token *oauth2.Token) (string, error) {
	// Make a request to the GitHub API with to get the user's emails
	client := p.config.Client(ctx, token)
	resp, err := client.Get(p.emailsUrl)

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
