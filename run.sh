#!/bin/sh
SAVED_GOPATH=$GOPATH
export GOPATH=$(pwd)

if [ "$1" = "build" ]; then
    go build -o bin/main src/main.go
elif [ "$1" = "fmt" ]; then
	gofmt -w src
	gofmt -w src/application
	gofmt -w src/bll
	gofmt -w src/bll/common
	gofmt -w src/bll/report
	gofmt -w src/bll/testutil
	gofmt -w src/dal
	gofmt -w src/util
elif [ "$1" = "lint" ]; then
	golint src
	golint src/application
	golint src/bll
	golint src/bll/common
	golint src/bll/report
	golint src/bll/testutil
	golint src/dal
	golint src/util
elif [ "$1" = "test" ]; then
    goconvey .
fi

export GOPATH=$SAVED_GOPATH
