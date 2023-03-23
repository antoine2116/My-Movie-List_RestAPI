package config

import (
	"fmt"

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
	Db string `mapstructure:"db"`
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

func Load(path string) (*Config, error) {
	viper.SetConfigName("local")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading configuration file : %s", err)
	}

	var c *Config

	err := viper.Unmarshal(&c)

	if err != nil {
		return nil, fmt.Errorf("unable to decode configuration file into struct: %s", err)
	}

	return c, nil
}
