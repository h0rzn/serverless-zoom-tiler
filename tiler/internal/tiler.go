package internal

import (
	"archive/zip"
	"bytes"
	"io"
	"log/slog"
	"strings"

	"github.com/lvisei/image2tiles"
)

type TilingOptions struct {
	Filename string
}

type TilingResult struct {
	ZipData io.Reader
}

func Tile(reader io.ReadCloser, opts TilingOptions) (TilingResult, error) {
	slog.Info("Tile", "filename", opts.Filename)
	defer reader.Close()

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	filename := opts.Filename
	if filename == "" {
		// TODO: find better default
		filename = "unkown-data"
	}

	writer, err := zw.Create(filename)
	if err != nil {
		return TilingResult{}, err
	}

	if _, err := io.Copy(writer, reader); err != nil {
		return TilingResult{}, err
	}

	w1, err := zw.Create("file_A")
	if err != nil {
		return TilingResult{}, err
	}
	r1 := strings.NewReader("content for file_A")
	if _, err := io.Copy(w1, r1); err != nil {
		return TilingResult{}, err
	}

	if err := zw.Close(); err != nil {
		return TilingResult{}, err
	}

	return TilingResult{
		ZipData: bytes.NewReader(buf.Bytes()),
	}, nil
}

func TileDisk(filename string) (outPath string, err error) {
	converter := image2tiles.NewConverter()
	err = converter.Prepare(filename, "#00000000")
	if err != nil {
		return "", err
	}

	out := filename + "_proc" + "/tile_%d_%d_%d.png"

	err = converter.Execute([2]int{512, 512}, true, 20, out)
	if err != nil {
		return "", err
	}

	return "", nil
}
