#!/bin/sh
SAVED_GOPATH=$GOPATH
export GOPATH=$(pwd)

if [ "$1" = "build" ]; then
    go build -o bin/main src/main.go
elif [ "$1" = "test" ]; then
    goconvey .
fi

export GOPATH=$SAVED_GOPATH
