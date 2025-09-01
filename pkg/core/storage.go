package core

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hasssanezzz/go-trgm/pkg/common"
)

type storageManager struct {
	file *os.File
	mu   sync.RWMutex
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
	// Read the block size and the block bytes under a read lock using ReadAt
	// so we don't change file offset and allow concurrent readers.
	blockSizeBuff := make([]byte, 4)
	s.mu.RLock()
	if _, err := s.file.ReadAt(blockSizeBuff, offset); err != nil {
		s.mu.RUnlock()
		return nil, err
	}

	blockSize := binary.LittleEndian.Uint32(blockSizeBuff)
	blockBytes := make([]byte, blockSize)
	if _, err := s.file.ReadAt(blockBytes, offset+4); err != nil {
		s.mu.RUnlock()
		return nil, err
	}
	s.mu.RUnlock()

	return append(blockSizeBuff, blockBytes...), nil
}

func (s *storageManager) readBlockAndDeserialize(offset int64) ([]common.DocumentID, error) {
	// Read raw block bytes under a read lock, then deserialize off-lock.
	// This avoids holding the lock during decompression/CPU work.
	blockSizeBuff := make([]byte, 4)
	s.mu.RLock()
	if _, err := s.file.ReadAt(blockSizeBuff, offset); err != nil {
		s.mu.RUnlock()
		return nil, err
	}

	blockSize := binary.LittleEndian.Uint32(blockSizeBuff)
	data := make([]byte, blockSize)
	if _, err := s.file.ReadAt(data, offset+4); err != nil {
		s.mu.RUnlock()
		return nil, err
	}
	s.mu.RUnlock()

	// Provide a reader that matches previous behavior (block size header + compressed data)
	r := bytes.NewReader(append(blockSizeBuff, data...))
	var block common.Block
	if err := block.Deserialize(r); err != nil {
		return nil, err
	}

	return block.Entries(), nil
}

func (s *storageManager) writeBlock(block *common.Block) (int64, error) {
	data := block.Serialize()

	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the block location
	offset, err := s.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Write block data
	if _, err := s.file.Write(data); err != nil {
		return 0, err
	}

	return offset, nil
}

func (s *storageManager) close() error {
	return s.file.Close()
}
