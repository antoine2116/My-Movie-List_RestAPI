package oauth

import (
	"apous-films-rest-api/config"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type UserEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(cfg config.GitHubConfig) GitHubProvider {
	return GitHubProvider{
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

func (p *GitHubProvider) GetUserEmail(ctx context.Context, code string) (string, error) {
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

	var emails []UserEmail

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
