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

    actions   = ["SNS:Publish"]
    resources = ["new-upload-topic-${var.region}-${var.stage}"]

    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"
      values   = [aws_s3_bucket.video_bucket.arn]
    }
  }
}

resource "aws_sns_topic" "new_upload_topic" {
  name   = "new-upload-topic-${var.region}-${var.stage}"
  policy = data.aws_iam_policy_document.new_upload_notification_policy.json
}

resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.video_bucket.id

  topic {
    topic_arn     = aws_sns_topic.new_upload_topic.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "raw"
  }
}