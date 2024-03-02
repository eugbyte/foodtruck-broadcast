
variable "lambda_function_name" {
  description = "function name of deployed lambda."
  type        = string
}

variable "lambda_arn" {
  description = "arn of lambda function."
  type        = string
}

variable "route_key" {
  description = "Route key for api gateway, e.g. 'GET /api/v1/subscription'"
  type        = string
}
