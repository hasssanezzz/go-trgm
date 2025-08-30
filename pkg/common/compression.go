//go:build !nocompress

package common

import (
	z "github.com/klauspost/compress/zstd"
)

var (
	compressor, errc   = z.NewWriter(nil, z.WithEncoderLevel(z.SpeedFastest))
	decompressor, errd = z.NewReader(nil)
)

func init() {
	if errc != nil {
		panic(errc)
	}
	if errd != nil {
		panic(errd)
	}
}

func compress(data []byte) ([]byte, error) {
	return compressor.EncodeAll(data, nil), nil
}

func decompress(data []byte) ([]byte, error) {
	return decompressor.DecodeAll(data, nil)
}
