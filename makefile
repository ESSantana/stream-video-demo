.PHONY: build clean

GO_COMPILE := env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w"

build:
	export GO111MODULE=on
	export CGO_ENABLE=0
	export CC=gcc

	${GO_COMPILE} -o bin/api/main cmd/serverless/api/main.go
	chmod +x bin/api/main
	mv bin/api/main bin/api/bootstrap
	cd bin/api/ && zip bootstrap.zip bootstrap

clean:
	rm -rf bin/