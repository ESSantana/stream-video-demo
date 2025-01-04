#BUILD GO APP
FROM golang:1.23-bookworm AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLE=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/main.go

# SETUP CONTAINER RELEASE
FROM gcr.io/distroless/base-debian12 AS release-stage
WORKDIR /
COPY --from=build-stage /server /server
EXPOSE 8080
USER root:root
ENTRYPOINT ["./server"]
