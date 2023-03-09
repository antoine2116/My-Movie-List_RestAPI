package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port   int    `mapstructure:"port"`
	Secret string `mapstructure:"secret"`
}

type DatabaseConfiguration struct {
	URI  string `mapstructure:"uri"`
	Dev  string `mapstructure:"dev"`
	Test string `mapstructure:"test"`
}

type GoogleOuathConfiguration struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type Configuration struct {
	Server   ServerConfiguration      `mapstructure:"server"`
	Database DatabaseConfiguration    `mapstructure:"database"`
	Google   GoogleOuathConfiguration `mapstructure:"google_oauth"`
}

func LoadConfiguration(path string) *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	var configuration *Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file : %s", err)
	}

	err := viper.Unmarshal(&configuration)

	if err != nil {
		log.Fatalf("Unable to decode configuration file into struct: %v", err)
	}

	return configuration
}
