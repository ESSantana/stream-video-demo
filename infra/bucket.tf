resource "aws_s3_bucket" "video_bucket" {
  bucket = "essantana-videos-${var.aws_region}-${var.stage}"
  
  tags = {
    Name        = "essantana-videos-${var.aws_region}-${var.stage}"
    Environment = var.stage
  }
}

data "aws_iam_policy_document" "new_upload_notification_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }

    actions   = [
      "sns:Publish"
    ]
    resources = ["arn:aws:sns:${var.aws_region}:${local.account_id}:new-upload-topic-${var.aws_region}-${var.stage}"]

    condition {
      test     = "ArnLike"
      variable = "AWS:SourceArn"
      values   = [aws_s3_bucket.video_bucket.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceAccount"
      values   = [local.account_id]
    }
  }
}

resource "aws_sns_topic" "new_upload_topic" {
  name   = "new-upload-topic-${var.aws_region}-${var.stage}"
  policy = data.aws_iam_policy_document.new_upload_notification_policy.json
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.video_bucket.id

  topic {
    topic_arn     = aws_sns_topic.new_upload_topic.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix   = "raw"
  }
}