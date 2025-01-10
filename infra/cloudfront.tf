resource "aws_cloudfront_distribution" "video_stream_distribution" {
  origin {
    domain_name              = aws_s3_bucket.video_stream.bucket_regional_domain_name
    origin_id                = aws_s3_bucket.video_stream.id
    origin_access_control_id = aws_cloudfront_origin_access_control.stream_video_distribution_oac.id
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Video Stream Files Distribution"

  default_cache_behavior {
    target_origin_id          = aws_s3_bucket.video_stream.id
    allowed_methods           = ["GET", "HEAD", "OPTIONS"]
    cached_methods            = ["GET", "HEAD"]
    compress                  = true
    cache_policy_id           = data.aws_cloudfront_cache_policy.cache_policy.id
    viewer_protocol_policy    = "allow-all"
    origin_request_policy_id  = data.aws_cloudfront_origin_request_policy.origin_request_policy.id
  }

  price_class = "PriceClass_200"

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags = {
    Name = "video_stream_distribution-${var.stage}"
    Environment = var.stage
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_ssm_parameter" "ssm_cloudfront_distribution" {
  name  = "/video-stream/cloudfront/distribution"
  type  = "String"
  value = aws_cloudfront_distribution.video_stream_distribution.domain_name
}

resource "aws_cloudfront_origin_access_control" "stream_video_distribution_oac" {
  name                              = "stream-video-distribution-oai"
  description                       = "OAC to access S3 stream video bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

data "aws_cloudfront_cache_policy" "cache_policy" {
  name = "Managed-CachingOptimized"
}

data "aws_cloudfront_origin_request_policy" "origin_request_policy" {
  name = "Managed-CORS-CustomOrigin"
}