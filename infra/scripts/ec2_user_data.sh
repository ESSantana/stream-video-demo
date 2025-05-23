#!bin/bash

sudo su
yum update -y
yum install -y docker
systemctl enable docker
systemctl start docker.service
usermod -a -G docker ec2-user

export VIDEO_BUCKET=${VIDEO_BUCKET}
export CLOUDFRONT_DIST=${CLOUDFRONT_DIST}
export STAGE=${STAGE}
export SERVER_PORT=${SERVER_PORT}

source /etc/bashrc

docker run -p 80:${SERVER_PORT} -d --name stream-video-demo \
  -e SERVER_PORT \
  -e VIDEO_BUCKET \
  -e CLOUDFRONT_DIST \
  -e STAGE \
  -v /tmp:/tmp:rw
emersonsantanadev/stream-video-demo:latest
