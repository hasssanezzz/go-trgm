package trgm

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Index struct {
	counter [TriMaxCount]uint32
	mapper  map[uint32]set
}

func newIndex() *Index {
	return &Index{
		counter: [TriMaxCount]uint32{},
		mapper:  map[uint32]set{},
	}
}

// FIXME: This function is a major indexing performance bottleneck.
// Gotta use a tree or something and eliminate the usage of maps eveywhere.
func (idx *Index) put(tri uint32, entry IndexEntry) {
	if idx.counter[tri] == 0 {
		idx.mapper[tri] = newSet()
		idx.mapper[tri].add(entry)
		idx.counter[tri]++
	}

	if !idx.mapper[tri].contains(entry) {
		idx.mapper[tri].add(entry)
		idx.counter[tri]++
	}
}

func (idx *Index) count(tri uint32) uint32 {
	return idx.counter[tri]
}

func (idx *Index) get(tri uint32) (uint32, set) {
	count := idx.counter[tri]
	if count == 0 {
		return 0, set{}
	}

	// DEBUG
	set := idx.mapper[tri]
	if set.size() != int(count) {
		panic("counter issue")
	}

	return count, set
}

func (idx *Index) clear(tri uint32) {
	idx.counter[tri] = 0
	idx.mapper[tri] = newSet()
}

func (idx *Index) serialize(tri uint32) []byte {
	if idx.counter[tri] == 0 {
		return nil
	}
	entries := idx.mapper[tri].toSlice()
	result := bytes.NewBuffer(nil)

	// Block size [uint32]
	blockSize := 4*2 + 6*len(entries)
	result.Write(binary.LittleEndian.AppendUint32(nil, uint32(blockSize)))

	// Trigram as int [uint32]
	result.Write(binary.LittleEndian.AppendUint32(nil, tri))

	// Count
	result.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(entries))))

	// Writing index entries 6 Byte/Entry
	for _, entry := range entries {
		result.Write(entry.Encode())
	}

	return result.Bytes()
}

func (idx *Index) decodeBlock(tri uint32, data []byte) ([]IndexEntry, error) {
	if triFromData := binary.LittleEndian.Uint32(data[:4]); tri != triFromData {
		return nil, fmt.Errorf("tri/block mismatch wanted %q got %q", intToTri(tri), intToTri(triFromData))
	}

	count := binary.LittleEndian.Uint32(data[4:8])
	if count == 0 {
		return nil, fmt.Errorf("data coruption, block entry count = 0")
	}

	// TODO skip deleted entries

	data = data[8:]
	entries := make([]IndexEntry, count)
	for i := range count {
		entryBuff := data[i*6 : i*6+6]
		entries[i] = IndexEntry{
			batchId: binary.LittleEndian.Uint16(entryBuff[:2]),
			offset:  binary.LittleEndian.Uint32(entryBuff[2:]),
		}
	}
	return entries, nil
}
