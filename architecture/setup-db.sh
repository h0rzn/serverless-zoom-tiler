#!/bin/sh

TABLE_NAME="table01"
REGION="eu-central-1"

echo "Setting up Database (DynamoDB)..."
# create table
awslocal dynamodb create-table \
    --table-name $TABLE_NAME \
    --key-schema AttributeName=id,KeyType=HASH \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --billing-mode PAY_PER_REQUEST \
    --region $REGION


# insert item
awslocal dynamodb put-item \
    --table-name $TABLE_NAME \
    --item '{"id":{"S":"foo"}}' \
    --region $REGION 

echo "Done."
