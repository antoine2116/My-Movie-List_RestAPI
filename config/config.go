package config

import (
	"log"

	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	Port   int
	Secret string
}

type DatabaseConfiguration struct {
	URI string
}

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

func LoadConfiguration(path string) *Configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
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
