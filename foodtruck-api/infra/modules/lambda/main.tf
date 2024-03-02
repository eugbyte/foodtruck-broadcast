// build the binary for the lambda function in a specified path
resource "null_resource" "fn_binary" {
  triggers = {
    src_path_fn     = "${var.src_path_fn}"
    binary_path_fn  = "${var.binary_path_fn}"
    archive_path_fn = "${var.archive_path_fn}"
  }

  provisioner "local-exec" {
    command = "env GOOS=linux go build -o ${var.binary_path_fn} -ldflags='-s -w' ${var.src_path_fn}"
  }
}

// zip the binary, as we can use only zip files to AWS lambda
data "archive_file" "lambda" {
  depends_on = [null_resource.fn_binary]

  type        = "zip"
  source_file = var.binary_path_fn
  output_path = var.archive_path_fn
  # Impt to set chmod to 0777 to grant permission to make the file executable, should terraform be deployed from a windows environment.
  output_file_mode = "0777"
}

// create the lambda function from zip file
resource "aws_lambda_function" "lambda" {
  depends_on = [data.archive_file.lambda]

  function_name = var.fn_name
  description   = "register web push subscription"
  role          = var.aws_iam_lambda_role_arn
  handler       = "main"
  memory_size   = 128
  timeout       = 5

  filename         = var.archive_path_fn
  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "go1.x"

  # on destroy hook
  provisioner "local-exec" {
    when    = destroy
    command = "rm -rf bin/subscription"
  }
}

// create log group in cloudwatch to gather logs of our lambda function
resource "aws_cloudwatch_log_group" "log_group" {
  name              = "/aws/lambda/${aws_lambda_function.lambda.function_name}"
  retention_in_days = 7
}


