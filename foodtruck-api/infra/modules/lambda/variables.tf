
variable "src_path_fn" {
  description = "Path to go source file."
  type        = string
}

variable "binary_path_fn" {
  description = "Path to go binary file."
  type        = string
}

variable "archive_path_fn" {
  description = "Path to go zipped binary file."
  type        = string
}

variable "fn_name" {
  description = "Name of lambda function"
  type        = string
}

variable "aws_iam_lambda_role_arn" {
  description = "arn of aws iam role for lambdas"
  type        = string
}
