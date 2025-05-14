#BUILD GO APP
FROM golang:1.24-bookworm AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLE=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server/main.go

# GET ffmpeg AND MAKE SH AVAILABLE FOR RELEASE STAGE
FROM busybox:1.37.0-uclibc AS busybox
WORKDIR /
RUN wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz && tar -xf ffmpeg-release-amd64-static.tar.xz
RUN mv ffmpeg-7.0.2-amd64-static/ ./ffmpeg

# SETUP CONTAINER RELEASE
FROM gcr.io/distroless/base-debian12 AS release-stage
WORKDIR /
COPY --from=busybox /bin/sh /bin/sh
COPY --from=busybox /bin/ln /bin/ln
COPY --from=busybox /ffmpeg /usr/local/bin/ffmpeg
COPY --from=build-stage /server /server
EXPOSE 8080
USER root:root
RUN ln -s /usr/local/bin/ffmpeg/ffmpeg /usr/bin/ffmpeg
ENTRYPOINT ["./server"]
