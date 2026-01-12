package internal

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	dlBasePath   = "/tmp/"
	inputPrefix  = "uploads/"
	outputPrefix = "processed/"
)

func ProcessFile(ctx context.Context, s3Client *s3.Client, bucketName string, objectKey string) error {
	fileName := filepath.Base(objectKey)
	fileReader, err := DownloadReader(ctx, s3Client, bucketName, objectKey)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	opts := TilingOptions{Filename: "tmp_" + objectKey}
	result, err := Tile(fileReader, opts)
	if err != nil {
		return err
	}

	processedKey := strings.Replace(objectKey, inputPrefix, outputPrefix, 1)
	if processedKey == objectKey {
		processedKey = outputPrefix + fileName + ".zip"
	}

	err = Upload(ctx, s3Client, bucketName, processedKey, result.ZipData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: USE ME
func Cleanup(localFile string) error {
	return os.Remove(localFile)
}
