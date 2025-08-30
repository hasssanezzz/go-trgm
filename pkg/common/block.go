package common

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/hasssanezzz/go-trgm/pkg/bitset"
)

type Block struct {
	tri     uint32
	bitset  *bitset.Bitset
	entries []DocumentID // doesn't need to be a set
}

func NewBlock(tri uint32, entries []DocumentID) *Block {
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

	var err error
	if data, err = decompress(data); err != nil {
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

	b.entries = make([]DocumentID, 0, int(count))
	for range int(count) {
		entry := [16]byte{}
		if _, err := r.Read(entry[:]); err != nil {
			return err
		}
		b.entries = append(b.entries, entry)
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
		result.Write(entry[:])
	}

	compressed, _ := compress(result.Bytes())

	// Block size [uint32]
	blockSize := len(compressed)

	return append(binary.LittleEndian.AppendUint32(nil, uint32(blockSize)), compressed...)
}

func (b *Block) Entries() []DocumentID {
	deletedCount := len(b.entries) - b.bitset.CountOnes()
	result := make([]DocumentID, 0, len(b.entries)-deletedCount)
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

func (b *Block) Delete(target DocumentID) bool {
	for i, entry := range b.entries {
		if target.Compare(entry) {
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
