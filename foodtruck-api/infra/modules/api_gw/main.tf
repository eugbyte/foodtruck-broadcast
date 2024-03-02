// API Gateway
// create an API gateway for HTTP API
resource "aws_apigatewayv2_api" "lambda_api" {
  name          = "${var.lambda_function_name}-api"
  protocol_type = "HTTP"
  description   = "Serverless API gateway for HTTP API and AWS Lambda function"

  cors_configuration {
    allow_headers = ["*"]
    allow_methods = [
      "GET",
    ]
    allow_origins = [
      "*" // NOTE: here we should provide a particular domain, but for the sake of simplicity we will use "*"
    ]
    expose_headers = []
    max_age        = 0
  }
}

// create logs for API GW
resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.lambda_api.name}"

  retention_in_days = 7
}

// create a stage for API GW
resource "aws_apigatewayv2_stage" "lambda_api" {
  api_id = aws_apigatewayv2_api.lambda_api.id

  name        = "dev"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw.arn

    format = jsonencode({
      requestId               = "$context.requestId"
      sourceIp                = "$context.identity.sourceIp"
      requestTime             = "$context.requestTime"
      protocol                = "$context.protocol"
      httpMethod              = "$context.httpMethod"
      resourcePath            = "$context.resourcePath"
      routeKey                = "$context.routeKey"
      status                  = "$context.status"
      responseLength          = "$context.responseLength"
      integrationErrorMessage = "$context.integrationErrorMessage"
      }
    )
  }
  depends_on = [aws_cloudwatch_log_group.api_gw]
}

// Link the lambda function to invoke when specific HTTP request is made via API GW
resource "aws_apigatewayv2_integration" "lambda" {
  api_id = aws_apigatewayv2_api.lambda_api.id

  integration_uri  = var.lambda_arn #aws_lambda_function.lambda.arn
  integration_type = "AWS_PROXY"
}

// specify route that will be used to invoke lambda function
resource "aws_apigatewayv2_route" "lambda" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = var.route_key
  target    = "integrations/${aws_apigatewayv2_integration.lambda.id}"
}

// provide permission for API GW to invoke lambda function
resource "aws_lambda_permission" "lambda" {
  statement_id  = "api-gateway-${var.lambda_function_name}"
  action        = "lambda:InvokeFunction"
  function_name = var.lambda_function_name # aws_lambda_function.lambda.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_apigatewayv2_api.lambda_api.execution_arn}/*/*"
}
