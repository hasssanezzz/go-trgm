package trgm

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hasssanezzz/trigram-index/pkg/common"
)

var Threshold uint32 = 5 * 1000

type Indexer interface {
	Index(string, common.IndexEntry) error
	Fetch(string) ([]common.IndexEntry, error)
	Display()
	Close() error
}

type TrigramIndexer struct {
	index   *Index
	storage *storageManager
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

	storage, err := newStorageManager(filepath.Join(dirPath, "trgm.idx"))
	if err != nil {
		return nil, err
	}
	ti.storage = storage

	return ti, nil
}

func (i *TrigramIndexer) Index(s string, entry common.IndexEntry) error {
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

func (i *TrigramIndexer) Fetch(pattern string) ([]common.IndexEntry, error) {
	trigrams, err := validateInputAndExtractTrigrams(pattern)
	if err != nil {
		return nil, err
	}

	results := common.NewSet()

	for _, tri := range trigrams {
		// Search in the in-memory index
		if count, set := i.index.get(tri); count > 0 {
			results.Union(set)
		}

		// Search in blocks
		if offsets, found := i.bindex.index[tri]; found {
			if err := i.searchInBlock(offsets, &results); err != nil {
				log.Printf("failed to search in block[%q]: %v", intToTri(tri), err)
				continue
			}
		}
	}

	return results.ToSlice(), nil
}

func (i *TrigramIndexer) Display() {}

func (i *TrigramIndexer) Close() error {
	if err := i.bindex.close(); err != nil {
		return err
	}
	return i.storage.close()
}

func (i *TrigramIndexer) thresholdCheck(trigrams []uint32) error {
	// Is threshold exceeded?
	for _, tri := range trigrams {
		// TODO: handle partial failures
		if i.index.count(tri) >= Threshold {

			// Serialize the block
			block := common.NewBlock(tri, i.index.mapper[tri].ToSlice())

			// log.Printf("Tri: %q exceeded it threshold %d, block size: %d", intToTri(tri), i.index.count(tri), len(blockBytes))

			// Write the block
			offset, err := i.storage.writeBlock(block)
			if err != nil {
				return err
			}

			// Register the block location
			if err := i.bindex.putBlock(tri, uint32(offset)); err != nil {
				return fmt.Errorf("failed to index the newly created block: %v", err)
			}

			// Clear past entries
			i.index.clear(tri)
		}
	}

	return nil
}

func (i *TrigramIndexer) searchInBlock(offsets []uint32, set *common.Set) error {
	for _, offset := range offsets {
		entries, err := i.storage.readBlockAndDeserialize(int64(offset))
		if err != nil {
			return err
		}

		set.UnionSlice(entries)
	}

	return nil
}
