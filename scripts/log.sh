#!/bin/sh

awslocal logs filter-log-events --log-group-name /aws/lambda/S3TilingHandler
