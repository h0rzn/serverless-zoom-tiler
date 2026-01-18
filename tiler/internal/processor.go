package internal

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"image"
	"log/slog"
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
	slog.Info("ProcessFile", "bucket", bucketName, "key", objectKey)

	filename := filepath.Base(objectKey)
	fileReader, err := DownloadReader(ctx, s3Client, bucketName, objectKey)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	opts := TilingOptions{Filename: filename}
	result, err := Tile(fileReader, opts)
	if err != nil {
		return err
	}

	processedKey := strings.Replace(objectKey, inputPrefix, outputPrefix, 1)
	if processedKey == objectKey {
		processedKey = outputPrefix + filename + ".zip"
	}

	err = Upload(ctx, s3Client, bucketName, processedKey, result.ZipData)
	if err != nil {
		return err
	}

	return nil
}

func ProcessFileEfficient(ctx context.Context, s3Client *s3.Client, bucketName string, objectKey string) error {
	slog.Info("Starting Efficient Processing", "bucket", bucketName, "key", objectKey)

	fileReader, err := DownloadReader(ctx, s3Client, bucketName, objectKey)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer fileReader.Close()

	img, _, err := image.Decode(fileReader)
	if err != nil {
		return fmt.Errorf("image decode failed: %w", err)
	}

	tileChan := make(chan TileData, 50)

	var tilingErr error
	go func() {
		tilingErr = TileImage(img, tileChan)
	}()

	zipBuf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuf)

	for tile := range tileChan {
		f, err := zipWriter.Create(tile.Name)
		if err != nil {
			return fmt.Errorf("could not create zip entry: %w", err)
		}
		_, err = f.Write(tile.Data)
		if err != nil {
			return fmt.Errorf("could not write to zip: %w", err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	if tilingErr != nil {
		return fmt.Errorf("tiling error: %w", tilingErr)
	}

	processedKey := "processed/" + filepath.Base(objectKey) + ".zip"
	r := bytes.NewReader(zipBuf.Bytes())
	return Upload(ctx, s3Client, bucketName, processedKey, r)
}
