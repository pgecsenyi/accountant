package main

import "fmr/application"

func main() {

	app := &application.Application{}
	app.Initialize()
	app.Execute()
}
