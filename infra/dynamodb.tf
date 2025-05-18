resource "aws_dynamodb_table" "video_stream_table" {
  name           = "video-stream-demo-${var.stage}"
  billing_mode   = "PROVISIONED"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "video_id"
  range_key      = "video_name"

  attribute {
    name = "video_id"
    type = "S"
  }

  attribute {
    name = "video_name"
    type = "S"
  }

  tags = {
    Name        = "video-stream-demo"
    Environment = var.stage
  }
}