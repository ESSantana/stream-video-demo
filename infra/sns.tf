
resource "aws_sns_topic" "upload_notification_topic" {
  name   = "upload-notification-${var.aws_region}-${var.stage}"
  policy = data.aws_iam_policy_document.upload_notification_topic_policy.json
}

data "aws_iam_policy_document" "upload_notification_topic_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["s3.amazonaws.com"]
    }

    actions = [
      "sns:Publish"
    ]
    resources = ["arn:aws:sns:${var.aws_region}:${local.account_id}:upload-notification-${var.aws_region}-${var.stage}"]

    condition {
      test     = "ArnLike"
      variable = "AWS:SourceArn"
      values   = [aws_s3_bucket.video_stream_demo.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceAccount"
      values   = [local.account_id]
    }
  }
}

resource "aws_sns_topic_subscription" "upload_notification_topic_server_subscription" {
  topic_arn = aws_sns_topic.upload_notification_topic.arn
  protocol  = "http"
  endpoint  = "http://${aws_instance.stream_video.public_dns}/video-processor"

  depends_on = [
    aws_instance.stream_video
  ]
}
