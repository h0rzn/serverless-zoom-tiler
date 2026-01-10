#!/bin/sh

FUNC_NAME="S3TilingHandler"
ZIP="tiler_code.zip"
EXEC_ROLE="arn:aws:iam::000000000000:role/lambda-ex"

awslocal lambda create-function \
    --function-name $FUNC_NAME \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://tiler_code.zip \
    --role $EXEC_ROLE


awslocal lambda add-permission \
    --function-name $FUNC_NAME \
    --statement-id s3-trigger-permission \
    --action "lambda:InvokeFunction" \
    --principal s3.amazonaws.com \
    --source-arn arn:aws:s3:::sample-bucket

sleep 2

awslocal s3api put-bucket-notification-configuration \
    --bucket sample-bucket \
    --notification-configuration file://architecture/notification_trigger.json


# show logs
#   awslocal logs describe-log-groups
#   awslocal logs filter-log-events --log-group-name /aws/lambda/S3TilingHandler
