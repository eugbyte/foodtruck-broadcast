terraform {
  required_version = ">= 1.4.6"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 4.16"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "aws" {
  region = "ap-southeast-1"

  default_tags {
    tags = {
      app = "foodtruck-webpush"
    }
  }
}

// for zipping and the lambda files
locals {
  src_path_fn     = "${path.root}/../cmd"
  binary_path_fn  = "${path.root}/../bin"
  archive_path_fn = "${path.root}/../bin"
}

// Allow lambda service to assume (use) the role with such policy.
// generates an IAM policy document in JSON format for use with resources that expect policy documents such as aws_iam_policy.
data "aws_iam_policy_document" "assume_lambda_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

resource "aws_sqs_queue" "geo_notify" {
  name                       = "foodtruck-geo-notify"
  message_retention_seconds  = 36000
  visibility_timeout_seconds = 30
  redrive_policy = jsonencode({
    "deadLetterTargetArn" = aws_sqs_queue.geo_notify_dlq.arn,
    "maxReceiveCount"     = 1
  })
}


resource "aws_sqs_queue" "geo_notify_dlq" {
  name                       = "foodtruck-geo-notify-DLQ"
  message_retention_seconds  = 36000
  visibility_timeout_seconds = 30
}

// create lambda role, that lambda function can assume (use)
resource "aws_iam_role" "lambda" {
  name               = "foodtruck-api-lambda-role"
  description        = "Role for lambda to assume lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_lambda_role.json
}


// Attach policies to a single lambda role - many to one

// grant lambda the permission to write to Cloudwatch logs
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role = aws_iam_role.lambda.id
  # AWSLambdaBasicExecutionRole is an AWS managed policy that allows your Lambda function to write to CloudWatch logs.
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

// attach sqs policy to lambda role
resource "aws_iam_role_policy_attachment" "sqs_policy" {
  role       = aws_iam_role.lambda.id
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

module "subscription_lambda" {
  source = "./../../../infra/modules/lambda"

  fn_name                 = "foodtruck-subscription"
  src_path_fn             = "${local.src_path_fn}/subscription"
  binary_path_fn          = "${local.binary_path_fn}/subscription/main"
  archive_path_fn         = "${local.archive_path_fn}/subscription/main.zip"
  aws_iam_lambda_role_arn = aws_iam_role.lambda.arn

  depends_on = [aws_iam_role.lambda, aws_iam_role_policy_attachment.lambda_policy]
}

module "subscription_api_gw" {
  source = "./../../../infra/modules/api_gw"

  lambda_function_name = module.subscription_lambda.lambda_function_name
  lambda_arn           = module.subscription_lambda.lambda_arn
  route_key            = "POST /api/v1/subscription"
}

module "producer_lambda" {
  source = "./../../../infra/modules/lambda"

  fn_name                 = "foodtruck-producer"
  src_path_fn             = "${local.src_path_fn}/producer"
  binary_path_fn          = "${local.binary_path_fn}/producer/main"
  archive_path_fn         = "${local.archive_path_fn}/producer/main.zip"
  aws_iam_lambda_role_arn = aws_iam_role.lambda.arn

  depends_on = [aws_iam_role.lambda, aws_iam_role_policy_attachment.lambda_policy, aws_iam_role_policy_attachment.sqs_policy]
}

module "producer_api_gw" {
  source = "./../../../infra/modules/api_gw"

  lambda_function_name = module.producer_lambda.lambda_function_name
  lambda_arn           = module.producer_lambda.lambda_arn
  route_key            = "POST /api/v1/notification/{geohash}"
}

module "consumer_lambda" {
  source = "./../../../infra/modules/lambda"

  fn_name                 = "foodtruck-consumer"
  src_path_fn             = "${local.src_path_fn}/consumer"
  binary_path_fn          = "${local.binary_path_fn}/consumer/main"
  archive_path_fn         = "${local.archive_path_fn}/consumer/main.zip"
  aws_iam_lambda_role_arn = aws_iam_role.lambda.arn

  depends_on = [aws_iam_role.lambda, aws_iam_role_policy_attachment.lambda_policy, aws_iam_role_policy_attachment.sqs_policy]
}

resource "aws_lambda_event_source_mapping" "consumer_sqs_mapping" {
  batch_size                         = 10 // https://docs.aws.amazon.com/lambda/latest/dg/invocation-eventsourcemapping.html
  maximum_batching_window_in_seconds = 1
  event_source_arn                   = aws_sqs_queue.geo_notify.arn
  enabled                            = true
  function_name                      = module.consumer_lambda.lambda_arn
}
