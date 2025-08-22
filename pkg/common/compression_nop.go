//go:build nocompress

package common

func compress(data []byte) ([]byte, error) {
	return data, nil
}

func decompress(data []byte) ([]byte, error) {
	return data, nil
}
