package internal

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"image"
	"log/slog"
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

func ProcessEfficient(ctx context.Context, s3Client *s3.Client, bucket, key string) error {
	localFile := dlBasePath + filepath.Base(key)
	Download(s3Client, ctx, bucket, key, localFile)

	outDir, err := TileDisk(localFile)
	if err != nil {
		return err
	}

	zipData, err := ZipDir(outDir)
	if err != nil {
		return nil
	}

	return Upload(ctx, s3Client, bucket, "processed/"+key+".zip", zipData)
}

func ProcessFileEfficient(ctx context.Context, s3Client *s3.Client, bucketName string, objectKey string) error {
	slog.Info("Starting Efficient Processing", "bucket", bucketName, "key", objectKey)

	// 1. S3 Download als Stream
	fileReader, err := DownloadReader(ctx, s3Client, bucketName, objectKey)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer fileReader.Close()

	// 2. Bild in den RAM dekodieren
	img, _, err := image.Decode(fileReader)
	if err != nil {
		return fmt.Errorf("image decode failed: %w", err)
	}

	// 3. Tiling-Prozess starten
	// Wir puffern den Channel leicht ab (z.B. 50), um Workern Spielraum zu geben
	tileChan := make(chan TileData, 50)

	// In einer Goroutine starten, damit wir im Hauptthread lesen können
	var tilingErr error
	go func() {
		tilingErr = TileImage(img, tileChan)
	}()

	// 4. Kacheln verarbeiten (z.B. in eine ZIP schreiben oder einzeln hochladen)
	// Da deine alten Funktionen eine ZIP hochladen, bauen wir hier eine im RAM
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

	// 5. Upload der fertigen ZIP
	processedKey := "processed/" + filepath.Base(objectKey) + ".zip"
	r := bytes.NewReader(zipBuf.Bytes())
	return Upload(ctx, s3Client, bucketName, processedKey, r)
}

// TODO: USE ME
func Cleanup(localFile string) error {
	return os.Remove(localFile)
}
