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
  iam_instance_profile    = aws_iam_instance_profile.server_instance_profile.id

  tags = {
    Name = "server-${var.stage}-${var.aws_region}"
  }
}

resource "aws_security_group" "server_security_group" {
  name        = "server-security-group-${var.stage}-${var.aws_region}"
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

data "aws_iam_policy_document" "server_assume_role_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "server_role" {
  name = "server-role-${var.stage}-${var.aws_region}"
  assume_role_policy = data.aws_iam_policy_document.server_assume_role_policy.json
}

data "aws_iam_policy_document" "server_s3_access_policy" {
  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:PutObject"
    ]

    resources = [aws_s3_bucket.video_bucket.arn]
  }
}

resource "aws_iam_policy" "server_s3_access_policy" {
  name   = "server-s3-access-policy-${var.stage}-${var.aws_region}"
  policy = data.aws_iam_policy_document.server_s3_access_policy.json
}

resource "aws_iam_role_policy_attachment" "server_role_attachment" {
  role = aws_iam_role.server_role.id
  policy_arn = aws_iam_policy.server_s3_access_policy.arn
}

resource "aws_iam_instance_profile" "server_instance_profile" {
  name = "server-instance-profile-${var.stage}-${var.aws_region}"
  role = aws_iam_role.server_role.id
}

resource "aws_sns_topic_subscription" "new_upload_topic_subscription" {
  topic_arn = aws_sns_topic.new_upload_topic.arn
  protocol  = "http"
  endpoint  = "http://${aws_instance.server.public_dns}/video-processor"
}