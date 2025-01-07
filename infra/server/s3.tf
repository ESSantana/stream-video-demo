resource "aws_s3_bucket" "video_stream" {
  bucket = "video-stream-${var.aws_region}-${var.stage}"
  
  tags = {
    Name        = "video-stream-${var.aws_region}-${var.stage}"
    Environment = var.stage
  }
}

resource "aws_s3_bucket_cors_configuration" "video_stream_cors_configuration" {
  bucket = aws_s3_bucket.video_stream.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "HEAD"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket_notification" "video_stream" {
  bucket = aws_s3_bucket.video_stream.id

  topic {
    topic_arn     = aws_sns_topic.upload_notification_topic.arn
    events        = ["s3:ObjectCreated:*"]
    filter_prefix   = "raw"
  }
}

resource "aws_s3_bucket_policy" "video_stream_bucket_policy" {
  bucket = aws_s3_bucket.video_stream.id
  policy = data.aws_iam_policy_document.video_stream_bucket_policy_document.json
}

data "aws_iam_policy_document" "video_stream_bucket_policy_document" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }
    effect = "Allow"
    actions = [
      "s3:GetObject",
    ]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.video_stream_distribution.arn]
    }
    resources = ["${aws_s3_bucket.video_stream.arn}/*"]
  }
}
