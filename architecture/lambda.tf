resource "aws_lambda_function" "tiling_handler" {
  function_name = "S3TilingHandler"
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  role          = "arn:aws:iam::000000000000:role/lambda-ex"

  filename      = "${path.module}/../tiler_code.zip"

  source_code_hash = filebase64sha256("${path.module}/../tiler_code.zip")
}

resource "aws_lambda_permission" "allow_s3" {
  statement_id  = "s3-trigger-permission"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.tiling_handler.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.s3.arn
}


resource "aws_s3_bucket_notification" "bucket_notification" {
  bucket = aws_s3_bucket.s3.id

  lambda_function {
    lambda_function_arn = aws_lambda_function.tiling_handler.arn
    events              = ["s3:ObjectCreated:*"]
    filter_prefix       = "uploads/"
  }

  depends_on = [aws_lambda_permission.allow_s3]
}
