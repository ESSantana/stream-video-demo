data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "dummy_code" {
    type = "zip"
    output_path = "${path.module}/dummy_code.zip"

    source {
        content = "Hello, world!"
        filename = "dummy_code.txt"
    }
}

resource "aws_lambda_function" "lambda" {
  filename      = "${data.archive_file.dummy_code.output_path}"
  function_name = "${var.function_name}-${var.stage}-${var.aws_region}"
  handler       = var.handler
  runtime       = var.runtime
  memory_size   = var.memory_size
  role          = aws_iam_role.iam_for_lambda.arn

  reserved_concurrent_executions = var.reserved_concurrent_executions

  environment {
    variables = var.environment_variables
  }

  tags          = var.tags
}