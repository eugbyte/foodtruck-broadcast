output "lambda_arn" {
  description = "arn of lambda function."
  value       = aws_lambda_function.lambda.arn
}

output "lambda_function_name" {
  description = "function name of deployed lambda."
  value       = aws_lambda_function.lambda.function_name
}
