package common

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/hasssanezzz/go-trgm/pkg/bitset"
	z "github.com/klauspost/compress/zstd"
)

type Block struct {
	tri     uint32
	bitset  *bitset.Bitset
	entries []IndexEntry // doesn't need to be a set
}

func NewBlock(tri uint32, entries []IndexEntry) *Block {
	return &Block{
		tri:     tri,
		bitset:  bitset.New(len(entries)),
		entries: entries,
	}
}

func (b *Block) Deserialize(r io.Reader) error {
	var blockSize uint32
	if err := binary.Read(r, binary.LittleEndian, &blockSize); err != nil {
		return err
	}

	data := make([]byte, blockSize)
	if _, err := r.Read(data); err != nil {
		return err
	}

	reader, err := z.NewReader(nil)
	if err != nil {
		return err
	}
	if data, err = reader.DecodeAll(data, nil); err != nil {
		return err
	}

	r = bytes.NewReader(data) // to minimize IO?

	if err := binary.Read(r, binary.LittleEndian, &b.tri); err != nil {
		return err
	}

	var count uint32
	if err := binary.Read(r, binary.LittleEndian, &count); err != nil {
		return err
	}

	bitsetBytes := make([]byte, bitset.CalcDataSize(int(count)))
	if _, err := r.Read(bitsetBytes); err != nil {
		return err
	}
	b.bitset = bitset.FromBytes(int(count), bitsetBytes)

	b.entries = make([]IndexEntry, 0, int(count))
	for range int(count) {
		entryBuff := make([]byte, 6)
		if _, err := r.Read(entryBuff); err != nil {
			return err
		}

		b.entries = append(b.entries, NewIndexEntry(
			binary.LittleEndian.Uint16(entryBuff[:2]),
			binary.LittleEndian.Uint32(entryBuff[2:])))
	}

	return nil
}

func (b *Block) Serialize() []byte {
	result := bytes.NewBuffer(nil)

	// Trigram as int [uint32]
	result.Write(binary.LittleEndian.AppendUint32(nil, b.tri))

	// Count [uint32]
	result.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(b.entries))))

	// Serialize the bitset []byte
	result.Write(b.bitset.Bytes())

	// Writing index entries 6 Byte/Entry
	for _, entry := range b.entries {
		result.Write(entry.Encode())
	}

	writer, _ := z.NewWriter(nil) // ignoring error here
	compressed := writer.EncodeAll(result.Bytes(), nil)

	// Block size [uint32]
	blockSize := len(compressed)

	return append(binary.LittleEndian.AppendUint32(nil, uint32(blockSize)), compressed...)
}

func (b *Block) Entries() []IndexEntry {
	deletedCount := len(b.entries) - b.bitset.CountOnes()
	result := make([]IndexEntry, 0, len(b.entries)-deletedCount)
	for i, entry := range b.entries {
		if b.bitset.Test(i) {
			continue
		}
		result = append(result, entry)
	}
	return result
}

func (b *Block) DeleteByIndex(index int) {
	b.bitset.Set(index)
}

func (b *Block) Delete(target IndexEntry) bool {
	for i, entry := range b.entries {
		if entry.batchId == target.batchId && entry.offset == target.offset {
			b.bitset.Set(i)
			return true
		}
	}
	return false
}

func BlockFromReader(r io.Reader) (*Block, error) {
	b := &Block{}
	if err := b.Deserialize(r); err != nil {
		return nil, err
	}
	return b, nil
}
