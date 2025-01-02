resource "aws_s3_bucket" "video_bucket" {
  bucket = "essantana-videos-${var.aws_region}-${var.stage}"

  tags = {
    Name        = "essantana-videos-${var.aws_region}-${var.stage}"
    Environment = var.stage
  }
}