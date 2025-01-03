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

resource "aws_iam_policy" "lambda_logging_policy" {
  name   = "${var.function_name}-${var.stage}-${var.aws_region}-logging-policy"
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        Action : [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Effect : "Allow",
        Resource : "arn:aws:logs:*:*:*"
      }
    ]
  })
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam-for-${var.function_name}-${var.stage}-${var.aws_region}"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "function_logging_policy_attachment" {
  role = aws_iam_role.iam_for_lambda.id
  policy_arn = aws_iam_policy.lambda_logging_policy.arn
}


resource "aws_cloudwatch_log_group" "lambda_log_group" {
  name = "/aws/lambda/${var.function_name}-${var.stage}-${var.aws_region}"

  # set one week for all lambda log groups 
  retention_in_days = 7 
  
  tags = var.tags
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
  timeout       = var.timeout
  architectures = var.architectures

  role          = aws_iam_role.iam_for_lambda.arn

  reserved_concurrent_executions = var.reserved_concurrent_executions

  depends_on = [aws_cloudwatch_log_group.lambda_log_group]

  environment {
    variables = var.environment_variables
  }

  tags          = var.tags
}
