
# 
# !!! The handler is the name of the executable for go1.x runtime!
resource "aws_lambda_function" "stores" {
  function_name    = "stores"
  filename         = "${path.module}/../services/openfga_http/src/functions/bin/handler.zip"
  handler          = "handler"
  source_code_hash = sha256(filebase64("${path.module}/../services/openfga_http/src/functions/bin/handler.zip"))
  role             = aws_iam_role.stores.arn
  runtime          = "provided.al2023"
  memory_size      = 128
  timeout          = 20
  # link to the layer that have the OpenFGA server running 
  # layers = ["arn:aws:lambda:eu-central-1:690777408331:layer:lambda-cache-layer:25"]
  layers = [ aws_lambda_layer_version.lambda_openfga_layer.arn ]
}


resource "aws_lambda_layer_version" "lambda_openfga_layer" {
  layer_name          = "lambda-openfga-layer"
  filename            = "${path.module}/../services/openfga_http/src/ext/bin/extension.zip"

  # Optionally, specify compatible runtimes for your layer
  compatible_runtimes = ["go1.x"]

  # You can also specify a description and license info
  description         = "OpenFGA server extension layer"
  license_info       = "MIT"

  # Terraform can automatically compute the source_code_hash for you, ensuring that
  # the layer version is updated when the ZIP file changes.
  source_code_hash    = filebase64sha256("${path.module}/../services/openfga_http/src/ext/bin/extension.zip")
}

# Create the role to give permissions to Lambda if any required 
resource "aws_iam_role" "stores" {
  name               = "stores"
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }
}
POLICY
}

resource "aws_iam_policy_attachment" "lambda_basic_execution" {
  name       = "lambda_basic_execution_attachment"
  roles      = [aws_iam_role.stores.name]
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# integrate with API Gateway
resource "aws_apigatewayv2_api" "lambda" {
  name          = "serverless_lambda_gw"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "lambda" {
  api_id = aws_apigatewayv2_api.lambda.id

  name        = "openfga"
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
}

resource "aws_apigatewayv2_integration" "openfga" {
  api_id = aws_apigatewayv2_api.lambda.id
  integration_uri    = aws_lambda_function.stores.invoke_arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  passthrough_behavior = "WHEN_NO_TEMPLATES"  # Ensures full request body forwarding
}

resource "aws_apigatewayv2_route" "openfga" {
  api_id = aws_apigatewayv2_api.lambda.id
  route_key = "ANY /{stores+}"
  target    = "integrations/${aws_apigatewayv2_integration.openfga.id}"
}

resource "aws_cloudwatch_log_group" "api_gw" {
  name = "/aws/api_gw/${aws_apigatewayv2_api.lambda.name}"
  retention_in_days = 5
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.stores.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.lambda.execution_arn}/*/*"
}


output "base_url" {
  description = "Base URL for API Gateway stage where OpenFGA is listening."
  value = aws_apigatewayv2_stage.lambda.invoke_url
}