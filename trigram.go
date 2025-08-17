package trgm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const TriMaxCount = 68 * 68 * 68

var Threshold uint32 = 5

type Indexer interface {
	Index(string, IndexEntry) error
	Fetch(string) ([]IndexEntry, error)
	Display()
	Close() error
}

type TrigramIndexer struct {
	index   *Index
	file    *os.File
	bindex  *blockIndex
	dirPath string
}

func NewTrigramIndexer(dirPath string) (Indexer, error) {
	ti := &TrigramIndexer{
		index:   newIndex(),
		dirPath: dirPath,
	}

	meta, err := newBlockIndex(dirPath)
	if err != nil {
		return nil, err
	}
	ti.bindex = meta

	file, err := os.OpenFile(filepath.Join(dirPath, "trgm.idx"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open index file: %v", err)
	}
	ti.file = file

	return ti, nil
}

func (i *TrigramIndexer) Index(s string, entry IndexEntry) error {
	trigrams, err := validateInputAndExtractTrigrams(s)
	if err != nil {
		return err
	}

	// Index
	for _, tri := range trigrams {
		i.index.put(tri, entry)
	}

	// Is threshold exceeded?
	for _, tri := range trigrams {
		// TODO: handle partial failures
		if i.index.count(tri) >= Threshold {
			log.Printf("Tri: %q exceeded it threshold %d", intToTri(tri), i.index.count(tri))

			// Seek to the end of the file before writing
			if _, err := i.file.Seek(0, io.SeekEnd); err != nil {
				return err
			}

			// Get the block location
			offset, err := i.file.Seek(0, io.SeekCurrent)
			if err != nil {
				return err
			}

			// Write the block
			blockBytes := i.index.serialize(tri)
			if _, err := i.file.Write(blockBytes); err != nil {
				return err
			}

			// Register the block location
			if err := i.bindex.putBlock(tri, uint32(offset)); err != nil {
				return fmt.Errorf("failed to index the newly created block: %v", err)
			}

			// Clear past entries
			i.index.mapper[tri] = make(map[IndexEntry]struct{})
			i.index.counter[tri] = 0
		}
	}

	return nil
}

func (i *TrigramIndexer) Fetch(pattern string) ([]IndexEntry, error) {
	trigrams, err := validateInputAndExtractTrigrams(pattern)
	if err != nil {
		return nil, err
	}

	results := map[IndexEntry]struct{}{}
	for _, tri := range trigrams {
		// Search in the in-memory index
		if mp, found := i.index.mapper[tri]; found {
			for entry, _ := range mp {
				results[entry] = struct{}{}
			}
		}

		// Search in blocks
		if offsets, found := i.bindex.index[tri]; found {
			entries, err := i.searchInBlock(tri, offsets)
			if err != nil {
				log.Printf("failed to search in block[%q]: %v", intToTri(tri), err)
				continue
			}

			// TODO: pass a pointer to the result map insteaf of allocating memory
			for _, entry := range entries {
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
	d, err := json.MarshalIndent(i.bindex.index, "", "    ")
	if err != nil {
		panic(err)
	}

	os.WriteFile("display.temp.txt", d, 0644)
}

func (i *TrigramIndexer) Close() error {
	if err := i.bindex.close(); err != nil {
		return err
	}
	return i.file.Close()
}

func (i *TrigramIndexer) searchInBlock(tri uint32, offsets []uint32) ([]IndexEntry, error) {
	// TODO: isn't block's size fixed?

	results := map[IndexEntry]struct{}{}

	for _, offset := range offsets {
		if _, err := i.file.Seek(int64(offset), io.SeekStart); err != nil {
			return nil, err
		}

		blockSizeBuff := make([]byte, 4)
		if _, err := i.file.Read(blockSizeBuff); err != nil {
			return nil, err
		}

		blockSize := binary.LittleEndian.Uint32(blockSizeBuff)
		block := make([]byte, blockSize)
		if _, err := i.file.Read(block); err != nil {
			return nil, err
		}

		entries, err := i.index.decodeBlock(tri, block)
		if err != nil {
			return nil, fmt.Errorf("searchInBlock failed to decode block: %v", err)
		}

		for _, entry := range entries {
			results[entry] = struct{}{}
		}
	}

	entries := make([]IndexEntry, 0, len(results))
	for entry, _ := range results {
		entries = append(entries, entry)
	}

	return entries, nil
}
