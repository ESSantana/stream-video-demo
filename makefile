.PHONY: build clean

GO_COMPILE := env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w"

build:
	export GO111MODULE=on
	export CGO_ENABLE=1

	${GO_COMPILE} -o bin/api/main cmd/serverless/api/main.go
	chmod +x bin/api/main
	mv bin/api/main bin/api/bootstrap
	cd bin/api/ && zip api.zip bootstrap

	${GO_COMPILE} -o bin/jobs/process-video/main cmd/serverless/jobs/video-processor/main.go
	chmod +x bin/jobs/process-video/main
	mv bin/jobs/process-video/main bin/jobs/process-video/bootstrap
	cd bin/jobs/process-video/ && zip process-video.zip bootstrap

clean:
	rm -rf bin/