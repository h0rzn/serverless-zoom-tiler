#!/bin/sh

make zip
./scripts/update-function.sh
awslocal s3 cp assets/grundriss.pdf s3://sample-bucket
awslocal logs filter-log-events --log-group-name /aws/lambda/S3TilingHandler

