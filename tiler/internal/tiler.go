package internal

import (
	"archive/zip"
	"bytes"
	"io"
)

type TilingOptions struct {
	Filename string
}

type TilingResult struct {
	ZipData io.Reader
}

func Tile(reader io.ReadCloser, opts TilingOptions) (TilingResult, error) {
	defer reader.Close()

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	filename := opts.Filename
	if filename == "" {
		filename = "data.bin"
	}

	f, err := zw.Create(filename)
	if err != nil {
		return TilingResult{}, err
	}

	if _, err := io.Copy(f, reader); err != nil {
		return TilingResult{}, err
	}

	if err := zw.Close(); err != nil {
		return TilingResult{}, err
	}

	return TilingResult{
		ZipData: bytes.NewReader(buf.Bytes()),
	}, nil
}
