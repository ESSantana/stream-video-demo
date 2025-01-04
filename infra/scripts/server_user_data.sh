#!bin/bash

sudo su
yum update -y
yum install -y docker
systemctl start docker.service
usermod -a -G docker ec2-user

docker run -p 80:8080 -d --name video-streaming-demo \
    -e SERVER_PORT="8080" \
    -e VIDEO_BUCKET="essantana-videos-sa-east-1-production" \
    emersonsantanadev/video-streaming-demo:latest 