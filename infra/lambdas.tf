module "stream_video_api" { 
    source = "./modules/lambda-function"

    function_name = "stream-video-api"
    stage         = var.stage
    handler       = "bin/api/bootstrap" 
    environment_variables = {
      VIDEO_BUCKET = aws_s3_bucket.video_bucket.id
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

  statement {
    effect    = "Allow"
    actions   = ["logs:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "stream_video_role_policy" {
  name    = "stream-video-role-policy"
  role    = module.stream_video_api.lambda_role_id
  policy  = data.aws_iam_policy_document.stream_video_policy_document.json
}


module "process_video_job" { 
    source = "./modules/lambda-function"

    function_name = "process-video"
    stage         = var.stage
    handler       = "bin/jobs/process-video/bootstrap" 
    environment_variables = {
      VIDEO_BUCKET = aws_s3_bucket.video_bucket.id
    }

    tags = {
      STAGE = var.stage
    }
}

data "aws_iam_policy_document" "process_video_job_policy_document" {
  statement {
    effect    = "Allow"
    actions   = ["s3:*"]
    resources = ["*"]
  }

  statement {
    effect    = "Allow"
    actions   = ["logs:*"]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "process_video_job_role_policy" {
  name    = "process-video-role-policy"
  role    = module.process_video_job.lambda_role_id
  policy  = data.aws_iam_policy_document.process_video_job_policy_document.json
}

resource "aws_lambda_permission" "allow_bucket" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = module.process_video_job.lambda_arn
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.video_bucket.arn
}

resource "aws_s3_bucket_notification" "new_video_upload_notification" {
  bucket = aws_s3_bucket.video_bucket.id

  lambda_function {
    lambda_function_arn = module.process_video_job.lambda_arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "raw/"
    filter_suffix       = ".mp4"
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}