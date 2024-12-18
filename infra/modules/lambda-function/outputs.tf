output "lambda_name" {
  value       = "${var.function_name}-${var.stage}-${var.aws_region}"
  description = "Lambda function name"
}

output "lambda_arn" {
  value       = aws_lambda_function.lambda.arn
  description = "Lambda function ARN"
}

output "lambda_invoke_arn" {
  value       = aws_lambda_function.lambda.invoke_arn
  description = "Lambda function invoke ARN"
}

output "lambda_role_arn" {
  value       = aws_iam_role.iam_for_lambda.arn
  description = "Lambda role ARN"
}

output "lambda_role_id" {
  value       = aws_iam_role.iam_for_lambda.id
  description = "Lambda role id"
}

output "lambda_version" {
  value       = aws_lambda_function.lambda.version
  description = "Lambda function version"
}