// create dynamo table
resource "aws_dynamodb_table" "subscription" {
  name         = "subscription"
  billing_mode = "PAY_PER_REQUEST"

  hash_key = "endpoint"

  attribute {
    name = "endpoint"
    type = "S"
  }
  attribute {
    name = "geohash"
    type = "S"
  }
  attribute {
    name = "lastSend"
    type = "N"
  }

  global_secondary_index {
    name            = "geohashIndex"
    hash_key        = "geohash"
    range_key       = "lastSend"
    write_capacity  = 10
    read_capacity   = 10
    projection_type = "ALL"
  }
}

data "aws_iam_policy_document" "allow_dynamodb_table_operations" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:*",
    ]

    resources = [
      aws_dynamodb_table.subscription.arn,
    ]
  }
}

resource "aws_iam_policy" "dynamodb_lambda_policy" {
  name        = "subscription_DynamoDBLambdaPolicy"
  description = "Policy for lambda to operate on dynamodb table"
  policy      = data.aws_iam_policy_document.allow_dynamodb_table_operations.json
}

resource "aws_iam_role_policy_attachment" "lambda_dynamodb_policy_attachment" {
  role       = aws_iam_role.lambda.id
  policy_arn = aws_iam_policy.dynamodb_lambda_policy.arn

  depends_on = [aws_iam_role.lambda]
}
