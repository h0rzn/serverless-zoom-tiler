#!/bin/sh

BUCKET_NAME="cc2526-tiling-storage"
# awslocal s3 cp assets/grundriss.pdf s3://$BUCKET_NAME/uploads/grundriss.pdf
awslocal s3 cp assets/gr.jpg s3://$BUCKET_NAME/uploads/gr2.png
aws s3 cp assets/gr.jpg s3://$BUCKET_NAME/uploads/gr2.png
