package internal

import "io"

type TilingOptions struct{}

type TilingResult struct {
	ZipData io.Reader
}

func Tile(opts TilingOptions) (TilingResult, error) {
	return TilingResult{}, nil
}
