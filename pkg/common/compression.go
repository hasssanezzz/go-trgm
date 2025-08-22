//go:build !nocompress

package common

import (
	z "github.com/klauspost/compress/zstd"
)

func compress(data []byte) ([]byte, error) {
	writer, err := z.NewWriter(nil)
	if err != nil {
		return nil, err
	}

	return writer.EncodeAll(data, nil), nil
}

func decompress(data []byte) ([]byte, error) {
	reader, err := z.NewReader(nil)
	if err != nil {
		return nil, err
	}

	return reader.DecodeAll(data, nil)
}
