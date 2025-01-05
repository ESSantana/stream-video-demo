data "aws_caller_identity" "current" {}

locals {
    account_id = data.aws_caller_identity.current.account_id
}

resource "aws_instance" "server" {
  ami                     = "ami-03c4a8310002221c7"
  instance_type           = "t2.micro"
  user_data = file("./scripts/server_user_data.sh")
  vpc_security_group_ids  = [aws_security_group.server_security_group.id]
  key_name                = "aws-emerson-sa-east-1"
}

resource "aws_security_group" "server_security_group" {
  name        = "server_security_group"
  description = "Server Security Group"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

   ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_sns_topic_subscription" "new_upload_topic_subscription" {
  topic_arn = aws_sns_topic.new_upload_topic.arn
  protocol  = "http"
  endpoint  = "http://${aws_instance.server.public_dns}/video-processor"
}