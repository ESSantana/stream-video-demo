module "stream_video_api" { 
    source = "./modules/lambda-function"

    function_name = "stream-video-api"
    stage         = var.stage
    handler       = "bin/api/bootstrap" 
    environment_variables = {
      # VIDEO_BUCKET = aws_s3_bucket.video_bucket.id
    }

    tags = {
      STAGE = var.stage
    }
}

data "aws_iam_policy_document" "stream_video_policy_document" {
  statement {
    effect    = "Allow"
    actions   = ["s3:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "stream_video_role_policy" {
  name    = "stream-video-role-policy"
  role    = module.stream_video_api.lambda_role_id
  policy  = data.aws_iam_policy_document.stream_video_policy_document.json
}

resource "aws_api_gateway_rest_api" "stream_video_api" {
  name = "stream_video_api-${var.stage}"
}

resource "aws_api_gateway_resource" "stream_video_resource" {
  parent_id   = aws_api_gateway_rest_api.stream_video_api.root_resource_id
  path_part   = "{proxy+}"
  rest_api_id = aws_api_gateway_rest_api.stream_video_api.id
}

resource "aws_api_gateway_method" "stream_video_method" {
  authorization = "NONE"
  http_method   = "ANY"
  resource_id   = aws_api_gateway_resource.stream_video_resource.id
  rest_api_id   = aws_api_gateway_rest_api.stream_video_api.id
}

resource "aws_api_gateway_integration" "stream_video_integration" {
  http_method             = aws_api_gateway_method.stream_video_method.http_method
  resource_id             = aws_api_gateway_method.stream_video_method.resource_id
  rest_api_id             = aws_api_gateway_rest_api.stream_video_api.id
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
  uri                     = module.stream_video_api.lambda_invoke_arn
}

resource "aws_lambda_permission" "stream_video_lambda_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.stream_video_api.lambda_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_api_gateway_rest_api.stream_video_api.execution_arn}/*/*"
}

resource "aws_api_gateway_deployment" "stream_video_deployment" {
  rest_api_id = aws_api_gateway_rest_api.stream_video_api.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.stream_video_resource.id,
      aws_api_gateway_method.stream_video_method.id,
      aws_api_gateway_integration.stream_video_integration.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "stream_video_stage" {
  deployment_id = aws_api_gateway_deployment.stream_video_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.stream_video_api.id
  stage_name    = var.stage
}