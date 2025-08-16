package trgm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

type blockIndex struct {
	blocks   map[uint32][]uint32
	wal      *_WAL
	filename string
}

func newBlockIndex(dirPath string) (*blockIndex, error) {
	m := &blockIndex{
		blocks:   map[uint32][]uint32{},
		filename: filepath.Join(dirPath, "block_index.bin"),
	}

	wal, err := newBlockIndexWAL(filepath.Join(dirPath, "block_index.wal"))
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL: %v", err)
	}
	m.wal = wal

	blocks, err := wal.read()
	if err != nil {
		return nil, fmt.Errorf("failed to read WAL: %v", err)
	}
	if len(blocks) > 0 {
		m.blocks = blocks
	}

	if _, err := os.Stat(m.filename); err != nil {
		if os.IsNotExist(err) {
			return m, nil
		}

		return nil, fmt.Errorf("failed to read index block file stats: %v", err)
	}

	if err := m.parse(); err != nil {
		return nil, err
	}

	return m, nil
}

func (b *blockIndex) parse() error {
	data, err := os.ReadFile(b.filename)
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)

	// Read number of trigrams
	var triCount uint32
	if err := binary.Read(r, binary.LittleEndian, &triCount); err != nil {
		return fmt.Errorf("failed to read trigram count: %w", err)
	}

	// Read each trigram's metadata
	for t := uint32(0); t < triCount; t++ {
		var tri uint32
		if err := binary.Read(r, binary.LittleEndian, &tri); err != nil {
			return fmt.Errorf("failed to read trigram id: %w", err)
		}

		var count uint32
		if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
			return fmt.Errorf("failed to read offset count: %w", err)
		}

		offsets := make([]uint32, count)
		for j := uint32(0); j < count; j++ {
			var o uint32
			if err := binary.Read(r, binary.LittleEndian, &o); err != nil {
				return fmt.Errorf("failed to read offset: %w", err)
			}
			offsets[j] = o
		}

		b.blocks[tri] = offsets
	}

	return nil
}

func (b *blockIndex) serialize() error {
	buffer := bytes.NewBuffer(nil)

	// Write number of trigrams
	buffer.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(b.blocks))))

	for tri, offsets := range b.blocks {
		buffer.Write(binary.LittleEndian.AppendUint32(nil, tri))
		buffer.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(offsets))))
		for _, offset := range offsets {
			buffer.Write(binary.LittleEndian.AppendUint32(nil, offset))
		}
	}

	if err := os.WriteFile(b.filename, buffer.Bytes(), 0644); err != nil {
		return err
	}

	// If the block index is successfully serialized, then we can safely clear the WAL
	if err := b.wal.clear(); err != nil {
		return err
	}

	return nil
}

func (b *blockIndex) putBlock(tri, offset uint32) error {
	b.blocks[tri] = append(b.blocks[tri], offset)
	return b.wal.append(tri, offset)
}

func (b *blockIndex) close() error {
	if err := b.serialize(); err != nil {
		return err
	}

	if err := b.wal.close(); err != nil {
		return err
	}

	return nil
}
