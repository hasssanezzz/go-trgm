package core

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hasssanezzz/go-trgm/pkg/common"
)

type storageManager struct {
	file *os.File
	mu   sync.Mutex
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

func (s *storageManager) readBlock(offset int64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.file.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, err
	}

	blockSizeBuff := make([]byte, 4)
	if _, err := s.file.Read(blockSizeBuff); err != nil {
		return nil, err
	}

	blockSize := binary.LittleEndian.Uint32(blockSizeBuff)
	blockBytes := make([]byte, blockSize)
	if _, err := s.file.Read(blockBytes); err != nil {
		return nil, err
	}

	return append(blockSizeBuff, blockBytes...), nil
}

func (s *storageManager) readBlockAndDeserialize(offset int64) ([]common.DocumentID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.file.Seek(int64(offset), io.SeekStart); err != nil {
		return nil, err
	}

	var block common.Block
	if err := block.Deserialize(s.file); err != nil {
		return nil, err
	}

	return block.Entries(), nil
}

func (s *storageManager) writeBlock(block *common.Block) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the block location
	offset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Write block data
	if _, err := s.file.Write(block.Serialize()); err != nil {
		return 0, nil
	}

	return offset, nil
}

func (s *storageManager) close() error {
	return s.file.Close()
}
