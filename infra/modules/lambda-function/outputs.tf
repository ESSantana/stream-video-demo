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

output "lambda_version" {
  value       = aws_lambda_function.lambda.version
  description = "Lambda function version"
}