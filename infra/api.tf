
resource "aws_api_gateway_rest_api" "stream_video_api" {
  name = "stream-video-api-${var.stage}"
}

resource "aws_api_gateway_resource" "stream_video_resource" {
  parent_id   = aws_api_gateway_rest_api.stream_video_api.root_resource_id
  path_part   = "stream"
  rest_api_id = aws_api_gateway_rest_api.stream_video_api.id
}

resource "aws_api_gateway_resource" "upload_video_resource" {
  parent_id   = aws_api_gateway_resource.stream_video_resource.id
  path_part   = "upload"
  rest_api_id = aws_api_gateway_rest_api.stream_video_api.id
}

resource "aws_api_gateway_method" "upload_video_method" {
  authorization = "NONE"
  http_method   = "POST"
  resource_id   = aws_api_gateway_resource.upload_video_resource.id
  rest_api_id   = aws_api_gateway_rest_api.stream_video_api.id
}

resource "aws_api_gateway_integration" "upload_video_integration" {
  http_method             = aws_api_gateway_method.upload_video_method.http_method
  resource_id             = aws_api_gateway_resource.upload_video_resource.id
  rest_api_id             = aws_api_gateway_rest_api.stream_video_api.id
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = module.upload_video.lambda_invoke_arn
}

resource "aws_api_gateway_deployment" "upload_video_deployment" {
  rest_api_id = aws_api_gateway_rest_api.stream_video_api.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.upload_video_resource.id,
      aws_api_gateway_method.upload_video_method.id,
      aws_api_gateway_integration.upload_video_integration.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "stream_video_stage" {
  deployment_id = aws_api_gateway_deployment.upload_video_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.stream_video_api.id
  stage_name    = var.stage
}