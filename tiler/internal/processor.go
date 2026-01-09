package internal

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const dlBasePath = "/tmp/"

func ProcessFile(ctx context.Context, s3Client *s3.Client, bucketName string, objectKey string) error {
	dlPath := dlBasePath + objectKey
	err := Download(s3Client, ctx, bucketName, objectKey, dlPath)
	if err != nil {
		return err
	}

	opts := TilingOptions{}
	result, err := Tile(opts)
	if err != nil {
		return err
	}

	processedKey := objectKey + ".zip"
	err = Upload(ctx, s3Client, bucketName, processedKey, result.ZipData)
	if err != nil {
		return err
	}

	return nil
}

func Cleanup() error {
	// TODO: remove old files etc.
	return nil
}
