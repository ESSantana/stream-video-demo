data "aws_caller_identity" "current" {}

locals {
    account_id = data.aws_caller_identity.current.account_id
}

resource "aws_instance" "stream_video" {
  ami                     = "ami-0b5a42ccb0a949cf1" # Replace with another id if you want use another AMI, are not in the sa-east-1 region or it changed
  instance_type           = "t2.micro"
  user_data = file("./scripts/ec2_user_data.sh")
  vpc_security_group_ids  = [
    aws_security_group.sg_web_access_stream_video.id,
    aws_security_group.sg_remove_access_stream_video.id,
  ]
  key_name                = "aws-emerson-sa-east-1" # Replace with another key
  iam_instance_profile    = aws_iam_instance_profile.stream_video_instance_profile.id

  tags = {
    Name = "stream-video-${var.aws_region}-${var.stage}"
  }
}

resource "aws_security_group" "sg_web_access_stream_video" {
  name        = "web-access-stream-video"
  description = "Security group to allow access to the web server from internet"

  ingress {
    from_port   = 80
    to_port     = 80
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

resource "aws_security_group" "sg_remove_access_stream_video" {
  name        = "remove-access-stream-video"
  description = "Security group to allow access to the server via SSH"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

data "aws_iam_policy_document" "stream_video_assume_role_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "stream_video_role" {
  name = "stream-video-role-${var.aws_region}-${var.stage}"
  assume_role_policy = data.aws_iam_policy_document.stream_video_assume_role_policy.json
}

data "aws_iam_policy_document" "stream_video_permissions_policy_document" {
  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:PutObject", 
      "s3:ListBucket"
    ]

    resources = [
      aws_s3_bucket.video_stream.arn,
      "${aws_s3_bucket.video_stream.arn}/*"
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "sns:ConfirmSubscription"
    ]

    resources = [aws_sns_topic.upload_notification_topic.arn]
  }

  statement {
    effect = "Allow"

    actions = [
      "ssm:GetParameter"
    ]

    resources = ["*"]
  }
}

resource "aws_iam_policy" "stream_video_permissions_policy" {
  name   = "stream-video-permissions-policy-${var.aws_region}-${var.stage}"
  policy = data.aws_iam_policy_document.stream_video_permissions_policy_document.json
}

resource "aws_iam_role_policy_attachment" "stream_video_role_attachment" {
  role = aws_iam_role.stream_video_role.id
  policy_arn = aws_iam_policy.stream_video_permissions_policy.arn
}

resource "aws_iam_instance_profile" "stream_video_instance_profile" {
  name = "stream-video-instance-profile-${var.aws_region}-${var.stage}"
  role = aws_iam_role.stream_video_role.id
}
