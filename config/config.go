package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port          int    `mapstructure:"port"`
	Secret        string `mapstructure:"secret"`
	TokenDuration int    `mapstructure:"token_duration"` // in hours
}

type ClientConfiguration struct {
	URI string `mapstructure:"uri"`
}

type DatabaseConfiguration struct {
	URI  string `mapstructure:"uri"`
	Dev  string `mapstructure:"dev"`
	Test string `mapstructure:"test"`
}

type GoogleOAuthConfiguration struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type GitHubOAuthConfiguration struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type Configuration struct {
	Server   ServerConfiguration      `mapstructure:"server"`
	Client   ClientConfiguration      `mapstructure:"client"`
	Database DatabaseConfiguration    `mapstructure:"database"`
	Google   GoogleOAuthConfiguration `mapstructure:"google_oauth"`
	GitHub   GitHubOAuthConfiguration `mapstructure:"github_oauth"`
}

var Config *Configuration

func LoadConfiguration(path string) *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file : %s", err)
	}

	err := viper.Unmarshal(&Config)

	if err != nil {
		log.Fatalf("Unable to decode configuration file into struct: %v", err)
	}

	return Config
}
