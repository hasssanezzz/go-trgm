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

	return i.thresholdCheck(trigrams)
}

func (i *TrigramIndexer) Fetch(pattern string) ([]IndexEntry, error) {
	trigrams, err := validateInputAndExtractTrigrams(pattern)
	if err != nil {
		return nil, err
	}

	results := newSet()

	for _, tri := range trigrams {
		// Search in the in-memory index
		if count, set := i.index.get(tri); count > 0 {
			results.union(set)
		}

		// Search in blocks
		if offsets, found := i.bindex.index[tri]; found {
			if err := i.searchInBlock(tri, offsets, &results); err != nil {
				log.Printf("failed to search in block[%q]: %v", intToTri(tri), err)
				continue
			}
		}
	}

	return results.toSlice(), nil
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

func (i *TrigramIndexer) thresholdCheck(trigrams []uint32) error {
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
			i.index.mapper[tri] = newSet()
			i.index.counter[tri] = 0
		}
	}

	return nil
}

func (i *TrigramIndexer) searchInBlock(tri uint32, offsets []uint32, set *set) error {
	// TODO: isn't block's size fixed?

	for _, offset := range offsets {
		if _, err := i.file.Seek(int64(offset), io.SeekStart); err != nil {
			return err
		}

		blockSizeBuff := make([]byte, 4)
		if _, err := i.file.Read(blockSizeBuff); err != nil {
			return err
		}

		blockSize := binary.LittleEndian.Uint32(blockSizeBuff)
		block := make([]byte, blockSize)
		if _, err := i.file.Read(block); err != nil {
			return err
		}

		entries, err := i.index.decodeBlock(tri, block)
		if err != nil {
			return fmt.Errorf("searchInBlock failed to decode block: %v", err)
		}

		set.unionSlice(entries)
	}

	return nil
}
