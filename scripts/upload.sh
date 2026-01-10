#!/bin/sh

BUCKET_NAME="sample-bucket"
awslocal s3 cp assets/grundriss.pdf s3://$BUCKET_NAME/grundriss.pdf
