package trgm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type block struct {
	tri     uint32
	entries []IndexEntry // doesn't need to be a set
}

func (b *block) read(r io.Reader) error {
	var blockSize uint32
	if err := binary.Read(r, binary.LittleEndian, &blockSize); err != nil {
		return err
	}

	data := make([]byte, blockSize)
	if _, err := r.Read(data); err != nil {
		return err
	}

	b.decode(data)
	return nil
}

func (b *block) encode() []byte {
	result := bytes.NewBuffer(nil)

	// Block size [uint32]
	blockSize := 4*2 + 6*len(b.entries)
	result.Write(binary.LittleEndian.AppendUint32(nil, uint32(blockSize)))

	// Trigram as int [uint32]
	result.Write(binary.LittleEndian.AppendUint32(nil, b.tri))

	// Count
	result.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(b.entries))))

	// Writing index entries 6 Byte/Entry
	for _, entry := range b.entries {
		result.Write(entry.Encode())
	}

	return result.Bytes()
}

func (b *block) decode(data []byte) error {
	b.tri = binary.LittleEndian.Uint32(data[:4])

	count := binary.LittleEndian.Uint32(data[4:8])
	if count == 0 {
		return fmt.Errorf("data coruption, block entry count = 0")
	}

	data = data[8:] // skip tri and count (8 bytes)
	entries := make([]IndexEntry, count)
	for i := range count {
		entryBuff := data[i*6 : i*6+6]
		entries[i] = IndexEntry{
			batchId: binary.LittleEndian.Uint16(entryBuff[:2]),
			offset:  binary.LittleEndian.Uint32(entryBuff[2:]),
		}
	}
	b.entries = entries
	return nil
}
