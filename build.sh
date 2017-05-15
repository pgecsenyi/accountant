#!/bin/sh
SAVED_GOPATH=$GOPATH
export GOPATH=$(pwd)
go build -o bin/main src/main.go
export GOPATH=$SAVED_GOPATH
