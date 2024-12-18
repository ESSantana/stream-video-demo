module "upload_video" { 
    source = "./modules/lambda-function"

    function_name = "upload-video"
    stage         = var.stage
    environment_variables = {
      VIDEO_BUCKET = "teste-lambda-vars-${var.stage}"
    }

    tags = {
      STAGE = var.stage
    }
}

data "aws_iam_policy_document" "upload_video_policy_document" {
  statement {
    effect    = "Allow"
    actions   = ["s3:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "upload_video_role_policy" {
  name    = "upload-video-role-policy"
  role    = module.upload_video.lambda_role_id
  policy  = data.aws_iam_policy_document.upload_video_policy_document.json
}

