@ECHO OFF

IF /i "%1" == "build" (
	cd src
	go build -o ../bin/main.exe
	cd ..
) ELSE IF /i "%1" == "fmt" (
	gofmt -w src
) ELSE IF /i "%1" == "lint" (
	golint src/...
) ELSE IF /i "%1" == "test" (
	goconvey .
)
