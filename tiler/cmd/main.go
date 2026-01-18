package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.rz.htw-berlin.de/cloud-computing/2025_26_WiSe/group13/serverless-zoom-tiler/tiler/internal"
)

var s3Client *s3.Client

func init() {
	client, err := internal.InitS3(context.Background())
	if err != nil {
		slog.Error("Failed to initialize S3 session", "error", err)
		// TODO: abort
	}
	s3Client = client
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	var info events.S3Event
	if err := json.Unmarshal(event, &info); err != nil {
		slog.Error("Failed to unmarshal event")
		return err
	}

	if len(info.Records) == 0 {
		slog.Warn("No Records in event")
		return nil
	}

	record := info.Records[0]
	bucketName := record.S3.Bucket.Name
	// TODO: URL decode! url.QueryUnescape
	objectKey := record.S3.Object.Key
	slog.Info("++DEBUG++", "bucket", bucketName, "key", objectKey)

	// err := internal.ProcessFile(ctx, s3Client, bucketName, objectKey)
	// err := internal.ProcessFileDisk(ctx, s3Client, bucketName, objectKey)
	err := internal.ProcessFileEfficient(ctx, s3Client, bucketName, objectKey)
	if err != nil {
		slog.Error("Failed to process file", "error", err)
	}
	// TODO: handle error
	// defer internal.Cleanup()

	return nil
}

func main() {
	lambda.Start(handleRequest)

	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	// exPath := filepath.Dir(ex)
	// fmt.Println(exPath)
	// converter := image2tiles.NewConverter()
	//
	// err = converter.Prepare("../../assets/gr1.jpg", "#00000000")
	// if err != nil {
	// 	slog.Error("error", err)
	// }
	//
	// err = converter.Execute([2]int{256, 256}, true, 75, "out/level-%d/%d-%d.jpg")
	// if err != nil {
	// 	slog.Error("error", err)
	// }
	//
	// fmt.Println(converter.MaxZoom, converter.TileSize)
	//
}
