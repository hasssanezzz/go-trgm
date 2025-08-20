package trgm

import (
	"fmt"
	"io"
	"os"
)

type storageManager struct {
	file *os.File
}

func newStorageManager(filepath string) (*storageManager, error) {
	s := &storageManager{}
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %v", err)
	}
	s.file = file

	return s, nil
}

func (s *storageManager) readBlock(offset int64) ([]IndexEntry, error) {
	if _, err := s.file.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, err
	}

	var block block
	if err := block.read(s.file); err != nil {
		return nil, err
	}

	return block.entries, nil
}

func (s *storageManager) writeBlock(data []byte) (int64, error) {
	// Get the block location
	offset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Write block data
	if _, err := s.file.Write(data); err != nil {
		return 0, nil
	}

	return offset, nil
}

func (s *storageManager) close() error {
	return s.file.Close()
}
