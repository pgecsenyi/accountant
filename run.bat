@ECHO OFF

SET SAVED_GOPATH=%GOPATH%
SET GOPATH=%CD%

IF /i "%1" == "build" (
	go build -o bin/main.exe src/main.go
) ELSE IF /i "%1" == "fmt" (
	gofmt -w src
	gofmt -w src/application
	gofmt -w src/bll
	gofmt -w src/bll/common
	gofmt -w src/bll/report
	gofmt -w src/bll/testutil
	gofmt -w src/dal
	gofmt -w src/util
) ELSE IF /i "%1" == "lint" (
	golint src
	golint src/application
	golint src/bll
	golint src/bll/common
	golint src/bll/report
	golint src/bll/testutil
	golint src/dal
	golint src/util
) ELSE IF /i "%1" == "test" (
	goconvey .
)

SET GOPATH=%SAVED_GOPATH%
