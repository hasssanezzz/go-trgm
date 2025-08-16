package trgm

import (
	"encoding/binary"
	"fmt"
	"os"
)

type _WAL struct {
	filepath string
	file     *os.File
}

func newBlockIndexWAL(filepath string) (*_WAL, error) {
	w := &_WAL{
		filepath: filepath,
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %v", filepath, err)
	}
	w.file = file

	return w, nil
}

func (w *_WAL) append(tri, offset uint32) error {
	data := append(binary.LittleEndian.AppendUint32(nil, tri), binary.LittleEndian.AppendUint32(nil, offset)...)
	if _, err := w.file.Write(data); err != nil {
		return fmt.Errorf("failed to append file %q: %v", w.filepath, err)
	}
	return nil
}

func (w *_WAL) read() (map[uint32][]uint32, error) {
	data, err := os.ReadFile(w.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %v", w.filepath, err)
	}

	result := map[uint32][]uint32{}
	for i := range len(data) / 8 {
		slice := data[i*8 : i*8+8]
		tri, offset := binary.LittleEndian.Uint32(slice[:4]), binary.LittleEndian.Uint32(slice[4:])

		if result[tri] == nil {
			result[tri] = []uint32{offset}
		} else {
			result[tri] = append(result[tri], offset)
		}
	}

	return result, nil
}

func (w *_WAL) clear() error {
	if err := w.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to clear WAL file %q: %v", w.filepath, err)
	}

	if _, err := w.file.Seek(0, 0); err != nil {
		// If seeking fails after truncation, the file state might be inconsistent
		// for the next write, so report the error.
		return fmt.Errorf("failed to seek to beginning of WAL file %q after clearing: %v", w.filepath, err)
	}

	return nil
}

func (w *_WAL) close() error {
	return w.file.Close()
}
