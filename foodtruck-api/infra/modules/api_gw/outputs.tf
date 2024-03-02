output "api_base_url" {
  description = "Base URL for API Gateway stage."
  value       = aws_apigatewayv2_stage.lambda_api.invoke_url
}

output "local_stack_base_url" {
  description = "local stack base url for development"
  value       = aws_apigatewayv2_api.lambda_api.id
}
