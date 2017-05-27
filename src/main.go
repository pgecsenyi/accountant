package main

import "application"

func main() {

	app := &application.Application{}
	app.Initialize()
	app.Execute()
}
