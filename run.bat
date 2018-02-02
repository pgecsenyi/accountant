@ECHO OFF

SET SAVED_GOPATH=%GOPATH%
SET GOPATH=%CD%

IF /i "%1" == "build" (
	go build -o bin/main.exe src/main.go
) ELSE IF /i "%1" == "test" (
	goconvey .
)

SET GOPATH=%SAVED_GOPATH%
