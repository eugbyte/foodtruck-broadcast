output "subscription_api_base_url" {
  description = "Base URL for subscription"
  value       = module.subscription_api_gw.api_base_url
}

output "subscription_local_stack_base_url" {
  description = "local stack base url for subscription"
  value       = "http://${module.subscription_api_gw.local_stack_base_url}.execute-api.localhost.localstack.cloud:4566/<PATH>"
}

output "producer_api_base_url" {
  description = "Base URL for API Gateway stage."
  value       = module.producer_api_gw.api_base_url
}

output "producer_local_stack_base_url" {
  description = "local stack base url for development"
  value       = "http://${module.producer_api_gw.local_stack_base_url}.execute-api.localhost.localstack.cloud:4566/<PATH>"
}
