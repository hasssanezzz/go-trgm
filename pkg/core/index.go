package core

import (
	"github.com/hasssanezzz/go-trgm/pkg/common"
)

type Mapper interface {
	Put(uint32, common.DocumentID)
	Get(uint32) (int, common.Set)
	Count(uint32) int
	Clear(uint32)
}

type Index struct {
	counter [TriMaxCount]int
	mapper  map[uint32]common.Set
}

func newIndex() Mapper {
	return &Index{
		counter: [TriMaxCount]int{},
		mapper:  make(map[uint32]common.Set, TriMaxCount/4),
	}
}

// FIXME: This function is a major indexing performance bottleneck.
// Gotta use a tree or something and eliminate the usage of maps eveywhere.
func (idx *Index) Put(tri uint32, docID common.DocumentID) {
	if idx.counter[tri] == 0 {
		idx.mapper[tri] = common.NewSet()
		idx.mapper[tri].Add(docID)
		idx.counter[tri] = 1
		return
	}

	set := idx.mapper[tri]
	oldSize := set.Size()
	set.Add(docID)
	if oldSize < set.Size() {
		idx.counter[tri]++
	}
}

func (idx *Index) Count(tri uint32) int {
	return idx.counter[tri]
}

func (idx *Index) Get(tri uint32) (int, common.Set) {
	count := idx.counter[tri]
	if count == 0 {
		return 0, common.Set{}
	}

	// DEBUG
	set := idx.mapper[tri]
	if set.Size() != int(count) {
		panic("counter issue")
	}

	return count, set
}

func (idx *Index) Clear(tri uint32) {
	idx.counter[tri] = 0
	idx.mapper[tri] = common.NewSet()
}
