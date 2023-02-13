package main

import (
	"apous-films-rest-api/app"
	"apous-films-rest-api/config"
)

func main() {
	// Configuration
	conf := config.LoadConfiguration()

	// Initialize application
	app := &app.App{}
	app.Initialize(conf)

	// Run
	app.Run(conf.Server.Port)
}
