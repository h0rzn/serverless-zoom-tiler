#!/bin/sh

FUNC_NAME="S3TilingHandler"
ZIP="tiler_code.zip"

awslocal lambda update-function-code \
    --function-name $FUNC_NAME \
    --zip-file fileb://tiler_code.zip

# show logs
#   awslocal logs describe-log-groups
#   awslocal logs filter-log-events --log-group-name /aws/lambda/S3TilingHandler
