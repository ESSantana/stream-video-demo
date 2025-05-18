#!bin/bash

sudo su
yum update -y
yum install -y docker
systemctl enable docker
systemctl start docker.service
usermod -a -G docker ec2-user

docker run -p 80:${SERVER_PORT} -d --name stream-video-demo \
    -e SERVER_PORT \
    -e VIDEO_BUCKET \
    -e CLOUDFRONT_DIST \
    emersonsantanadev/stream-video-demo:latest
 