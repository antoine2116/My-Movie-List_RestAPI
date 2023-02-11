package main

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/pkg/app"
	"log"

	"github.com/spf13/viper"
)

func main() {
	// Configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../../")
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration file : %s", err)
	}

	err := viper.Unmarshal(&configuration)

	if err != nil {
		log.Fatalf("Unable to decode configuration file into struct: %v", err)
	}

	// Initialize application
	app := &app.App{}
	app.Initialize(configuration)
	app.Run(configuration.Server.Port)
}
