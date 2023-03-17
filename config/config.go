package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port          int    `mapstructure:"port"`
	Secret        string `mapstructure:"secret"`
	TokenDuration int    `mapstructure:"token_duration"` // in hours
}

type ClientConfig struct {
	URI string `mapstructure:"uri"`
}

type DatabaseConfig struct {
	URI  string `mapstructure:"uri"`
	Dev  string `mapstructure:"dev"`
	Test string `mapstructure:"test"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type GitHubConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Client   ClientConfig   `mapstructure:"client"`
	Database DatabaseConfig `mapstructure:"database"`
	Google   GoogleConfig   `mapstructure:"google_oauth"`
	GitHub   GitHubConfig   `mapstructure:"github_oauth"`
}

func Load(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file : %s", err)
	}

	var c *Config

	err := viper.Unmarshal(&c)

	if err != nil {
		log.Fatalf("Unable to decode configuration file into struct: %v", err)
	}

	return c
}
