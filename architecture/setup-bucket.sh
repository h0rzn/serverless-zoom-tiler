#!/bin/sh
# Creates a new local S3 bucket,
# downloads a pdf file from web
# and uploads it the bucket.

BUCKET_NAME="sample-bucket"
ASSETS="./assets/"

echo "Setting up Storage (S3 Bucket)..."

# create bucket
awslocal s3api create-bucket --bucket $BUCKET_NAME

# download example pdf
curl https://www.grub-bau.de/fileadmin/grub-bau.de/download/Grundriss%20Erdgeschoss.pdf -o "$ASSETS/grundriss.pdf"

# upload example pdf
awslocal s3api put-object \
  --bucket sample-bucket \
  --key grundriss.pdf \
  --body "$ASSETS/grundriss.pdf"

echo "Done."
