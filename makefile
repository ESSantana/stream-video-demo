.PHONY: build clean

GO_COMPILE := env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w"

build:
	export GO111MODULE=on
	export CGO_ENABLE=0
	export CC=gcc

	${GO_COMPILE} -o bin/video-processor/main cmd/serverless/video-processor/main.go
	chmod +x bin/video-processor/main
	mv bin/video-processor/main bin/video-processor/bootstrap
	cd bin/video-processor/ && powershell Compress-Archive bootstrap video-processor.zip 

	${GO_COMPILE} -o bin/video-upload/main cmd/serverless/video-upload/main.go
	chmod +x bin/video-upload/main
	mv bin/video-upload/main bin/video-upload/bootstrap
	cd bin/video-upload/ && powershell Compress-Archive bootstrap video-upload.zip 

clean:
	rm -rf bin/

deploy: build clean