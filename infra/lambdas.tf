module "stream-video-api" { 
    source = "./modules/lambda-function"

    function_name = "stream-video-api"
    stage         = var.stage
    handler       = "bin/api/boostrap" 
    environment_variables = {
      VIDEO_BUCKET = "teste-lambda-vars-${var.stage}"
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
  role    = module.stream-video-api.lambda_role_id
  policy  = data.aws_iam_policy_document.stream_video_policy_document.json
}

