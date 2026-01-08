package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func init() {
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {
	// 	log.Fatalf("unable to load SDK config, %v", err)
	// }
	//
	// s3Client = s3.NewFromConfig(cfg)
}

func handleRequest(ctx context.Context, event json.RawMessage) error {
	log.Println("handleRequest...")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}
