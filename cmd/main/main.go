package main

import "apous-films-rest-api/pkg/app"

func main() {
	app := &app.App{}
	app.Initialize()
	app.Run()
}
