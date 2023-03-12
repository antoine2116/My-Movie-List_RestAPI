package oauth

import (
	"apous-films-rest-api/config"
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type GoogleUser struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Verified   bool   `json:"verified_email"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

type GoogleProvider struct {
	config *oauth2.Config
}

func NewGoogleProvider() *GoogleProvider {
	p := &GoogleProvider{}
	p.Init()
	return p
}

func (p *GoogleProvider) Init() {

	p.config = &oauth2.Config{
		ClientID:     config.Config.Google.ClientID,
		ClientSecret: config.Config.Google.ClientSecret,
		RedirectURL:  config.Config.Google.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}
}

func (p *GoogleProvider) GetGoogleUser(code string, user *GoogleUser) error {
	token, err := p.config.Exchange(context.TODO(), code, oauth2.AccessTypeOffline)

	if err != nil {
		return err
	}

	client := p.config.Client(context.TODO(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return err
	}

	return nil
}
