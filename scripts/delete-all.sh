#!/bin/sh

echo "Running clean up..."

for bucket in $(awslocal s3api list-buckets --query "Buckets[].Name" --output text); do
    echo "Delete Bucket: $bucket"
    awslocal s3 rb "s3://$bucket" --force
done

for table in $(awslocal dynamodb list-tables --query "TableNames[]" --output text); do
    echo "Delete Table: $table"
    awslocal dynamodb delete-table --table-name "$table"
done

echo "Done."
