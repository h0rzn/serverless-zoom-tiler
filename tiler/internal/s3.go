package internal

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// InitS3 initializes the global s3Client using the default AWS configuration.
func InitS3(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("unable to load SDK config: %v\n", err)
		return nil, err
	}

	return s3.NewFromConfig(cfg), nil
}

// Download retrieves an object from the S3 bucket and saves it
// to the local file system using the given filename.
// It requires a valid S3 client, a context for cancellation, and the bucket/key identifiers.
func Download(client *s3.Client, ctx context.Context, bucketName string, objectKey string, filename string) error {
	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})
	if err != nil {
		log.Printf("Failed to get object: %s: %v\n", objectKey, err)
		return err
	}

	defer output.Body.Close()
	err = writeLocalFile(filename, output.Body)
	return err
}

// UploadByFilename reads a file from the local path and uploads it to the given S3 bucket.
// Existence of the file inside the bucket is *not* explicitly checked afterwards.
func UploadByFilename(client *s3.Client, ctx context.Context, bucketName string, objectKey string, filename string) error {
	file, err := readLocalFile(filename)
	if err != nil {
		log.Printf("Failed to read local file %s: %v\n", filename, err)
		return err
	}
	defer file.Close()

	return Upload(ctx, client, bucketName, objectKey, file)
}

// Upload reads a file from io.Reader and uploads it to the given S3 bucket.
// Existence of the file inside the bucket is *not* explicitly checked afterwards.
func Upload(ctx context.Context, client *s3.Client, bucketName string, objectKey string, data io.Reader) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
		Body:   data,
	})
	if err != nil {
		log.Printf("Failed to put object %s: %v\n", objectKey, err)
		return err
	}

	return nil
}

// writeLocalFile is a helper that creates a local file and writes the contents
// of an io.ReadCloser into it.
func writeLocalFile(filename string, data io.ReadCloser) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	if err != nil {
		return err
	}
	return nil
}

// readLocalFile opens a local file for reading and returns the file pointer.
// The caller has to close the file.
func readLocalFile(filename string) (*os.File, error) {
	return os.Open(filename)
}
