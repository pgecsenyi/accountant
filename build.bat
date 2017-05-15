SET SAVED_GOPATH=%GOPATH%
SET GOPATH=%CD%
go build -o bin/main.exe src/main.go
SET GOPATH=%SAVED_GOPATH%
