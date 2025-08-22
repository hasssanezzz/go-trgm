package trgm

import (
	"github.com/hasssanezzz/go-trgm/pkg/common"
)

type Index struct {
	counter [TriMaxCount]uint32
	mapper  map[uint32]common.Set
}

func newIndex() *Index {
	return &Index{
		counter: [TriMaxCount]uint32{},
		mapper:  map[uint32]common.Set{},
	}
}

// FIXME: This function is a major indexing performance bottleneck.
// Gotta use a tree or something and eliminate the usage of maps eveywhere.
func (idx *Index) put(tri uint32, entry common.IndexEntry) {
	if idx.counter[tri] == 0 {
		idx.mapper[tri] = common.NewSet()
		idx.mapper[tri].Add(entry)
		idx.counter[tri]++
	}

	if !idx.mapper[tri].Contains(entry) {
		idx.mapper[tri].Add(entry)
		idx.counter[tri]++
	}
}

func (idx *Index) count(tri uint32) uint32 {
	return idx.counter[tri]
}

func (idx *Index) get(tri uint32) (uint32, common.Set) {
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

func (idx *Index) clear(tri uint32) {
	idx.counter[tri] = 0
	idx.mapper[tri] = common.NewSet()
}
