#BUILD GO APP
FROM golang:1.23-bookworm AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN GOARCH=amd64 CGO_ENABLE=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/main.go

# SETUP CONTAINER RELEASE
FROM scratch AS release-stage
WORKDIR /app
COPY --from=build-stage /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
