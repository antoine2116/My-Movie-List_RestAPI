package oauth

import (
	"apous-films-rest-api/config"
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

type GitHubUser struct {
	Email string `json:"email"`
}

type UserEmail struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider() *GitHubProvider {
	p := &GitHubProvider{}
	p.Init()
	return p
}

func (p *GitHubProvider) Init() {

	p.config = &oauth2.Config{
		ClientID:     config.Config.GitHub.ClientID,
		ClientSecret: config.Config.GitHub.ClientSecret,
		RedirectURL:  config.Config.GitHub.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

func (p *GitHubProvider) GetGitHubUser(code string, user *GitHubUser) error {
	// Exchange will do the handshake to retrieve the initial access token.
	token, err := p.config.Exchange(context.TODO(), code)

	if err != nil {
		return err
	}
	fmt.Println(token)
	// Make a request to the GitHub API with to get the user's profile.
	client := p.config.Client(context.TODO(), token)
	resp, err := client.Get("https://api.github.com/user/emails")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var emails []UserEmail

	err = json.NewDecoder(resp.Body).Decode(&emails)

	if err != nil {
		return err
	}

	for _, email := range emails {
		if email.Primary {
			user.Email = email.Email
			break
		}
	}

	return nil
}
