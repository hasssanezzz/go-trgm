package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strings"
)

const Threshold = 100

type Indexer interface {
	Index(string, IndexEntry) error
	Fetch(string) ([]IndexEntry, error)
	Display()
}

type IndexEntry struct {
	BatchId uint16
	Offset  uint32
}

type Index map[uint32]map[IndexEntry]struct{}

func (idx Index) serialize() []byte {
	// block size [uint32]
	// trigram [uint32]
	// number of entries [uint32]
	// list of entries[uint16 uint32][uint16 uint32][uint16 uint32]...
	result := bytes.NewBuffer(nil)

	for tri, entryMap := range idx {
		blockSize := 4*3 + 6*len(entryMap)
		result.Write(binary.LittleEndian.AppendUint32(nil, uint32(blockSize)))
		result.Write(binary.LittleEndian.AppendUint32(nil, tri))
		result.Write(binary.LittleEndian.AppendUint32(nil, uint32(len(entryMap))))

		for entry, _ := range idx[tri] {
			result.Write(binary.LittleEndian.AppendUint16(nil, entry.BatchId))
			result.Write(binary.LittleEndian.AppendUint32(nil, entry.Offset))
		}
	}

	return result.Bytes()
}

func trigramToInt(trigram string) uint32 {
	c0 := uint32(trigram[0] - 32)
	c1 := uint32(trigram[1] - 32)
	c2 := uint32(trigram[2] - 32)
	return c0*8836 + c1*94 + c2 // 8836 = 94*94
}

func intToTrigram(id uint32) string {
	c0 := id / 8836
	remainder := id % 8836
	c1 := remainder / 94
	c2 := remainder % 94
	return string([]byte{
		byte(c0 + 32),
		byte(c1 + 32),
		byte(c2 + 32),
	})
}

func extractTrigrams(s string) []uint32 {
	s = strings.ToLower(s)
	results := make([]uint32, len(s)-2)
	for i := range len(s) - 2 {
		tri := s[i : i+3]
		results[i] = trigramToInt(tri)
	}
	return results
}

type TrigramIndexer struct {
	index Index
}

func NewTrigramIndexer() Indexer {
	ti := &TrigramIndexer{
		index: map[uint32]map[IndexEntry]struct{}{},
	}

	return ti
}

func (i *TrigramIndexer) Index(s string, entry IndexEntry) error {
	trigrams := extractTrigrams(s)
	for _, tri := range trigrams {
		i.index[tri][entry] = struct{}{}
	}

	if len(i.index) > Threshold {

	}

	return nil
}

func (i *TrigramIndexer) Fetch(pattern string) ([]IndexEntry, error) {
	trigrams := extractTrigrams(pattern)

	results := map[IndexEntry]struct{}{}
	for _, tri := range trigrams {
		if mp, found := i.index[tri]; found {
			for entry, _ := range mp {
				results[entry] = struct{}{}
			}
		}
	}

	entries := make([]IndexEntry, 0, len(results))
	for entry, _ := range results {
		entries = append(entries, entry)
	}

	return entries, nil
}

func (i *TrigramIndexer) Display() {
	d, err := json.MarshalIndent(i.index, "", "    ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(d))
}
