#!bin/bash

sudo su
yum update -y
yum install -y docker
systemctl enable docker
systemctl start docker.service
usermod -a -G docker ec2-user

export VIDEO_BUCKET=$(aws ssm get-parameter --name "/video-stream/s3/bucket-name" --query "Parameter.Value" --output text)
export CLOUDFRONT_DIST=$(aws ssm get-parameter --name "/video-stream/cloudfront/distribution" --query "Parameter.Value" --output text)
source /etc/bashrc

docker run -p 80:8080 -d --name stream-video-demo \
    -e SERVER_PORT="8080" \
    -e VIDEO_BUCKET \
    -e CLOUDFRONT_DIST \
    emersonsantanadev/stream-video-demo:latest
