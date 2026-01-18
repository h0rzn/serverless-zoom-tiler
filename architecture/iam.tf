data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

# Role: Lambda executioner role
resource "aws_iam_role" "lambda_exec_role" {
  name               = "tiling_handler_execution_role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json
}



# Policy: Access to S3 Bucket (uploads/ and processed/)
resource "aws_iam_policy" "lambda_s3_tiling_policy" {
  name        = "LambdaS3TilingAccess"
  description = "Allows Lambda to read uploads/ and write to processed/"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = ["${aws_s3_bucket.s3.arn}/uploads/*"]
      },
      {
        Effect   = "Allow"
        Action   = ["s3:PutObject"]
        Resource = ["${aws_s3_bucket.s3.arn}/processed/*"]
      }
    ]
  })
}

# Policy: attach to role
resource "aws_iam_role_policy_attachment" "lambda_s3_attach" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_s3_tiling_policy.arn
}

# Group: developer group
resource "aws_iam_group" "developers" {
  name = "tiling-project-developers"
}

# Policy: Allow developer to see (+upload) everything in Bucket
resource "aws_iam_group_policy_attachment" "dev_s3_access" {
  group      = aws_iam_group.developers.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

resource "aws_iam_user" "dev_user" {
  name = "theo"
}

# Group: add users to group
resource "aws_iam_user_group_membership" "add_to_group" {
  user   = aws_iam_user.dev_user.name
  groups = [aws_iam_group.developers.name]
}
