#!/bin/sh

if [ "$1" = "build" ]; then
	cd src
	go build -o bin/main
	cd ..
elif [ "$1" = "fmt" ]; then
	gofmt -w src
elif [ "$1" = "lint" ]; then
	golint src/...
elif [ "$1" = "test" ]; then
	goconvey .
fi
