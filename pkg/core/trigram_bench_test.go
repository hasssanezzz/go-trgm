package core

import (
	"math/rand"
	"os"
	"testing"

	"github.com/hasssanezzz/go-trgm/pkg/common"
)

func makeDocs(n int) []common.DocumentID {
	docs := make([]common.DocumentID, n)
	for i := 0; i < n; i++ {
		var id common.DocumentID
		id[0] = byte(i)
		id[1] = byte(rand.Intn(256))
		docs[i] = id
	}
	return docs
}

func BenchmarkIndexPut(b *testing.B) {
	b.ReportAllocs()
	docs := makeDocs(1024)
	idx := newIndex()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Put(uint32(i%1000), docs[i%len(docs)])
	}
}

func BenchmarkStorage_WriteThenReadParallel(b *testing.B) {
	b.ReportAllocs()
	// Create a temp file for storage manager
	tmp, err := os.CreateTemp("", "trgm_bench_*.idx")
	if err != nil {
		b.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	sm, err := newStorageManager(tmp.Name())
	if err != nil {
		b.Fatalf("failed to create storage manager: %v", err)
	}
	defer sm.close()

	// create a block with many entries
	nEntries := 1000
	docs := makeDocs(nEntries)
	block := common.NewBlock(12345, docs)

	// Pre-write a number of blocks and collect offsets
	nPre := 256
	offsets := make([]int64, nPre)
	for i := 0; i < nPre; i++ {
		off, err := sm.writeBlock(block)
		if err != nil {
			b.Fatalf("writeBlock failed: %v", err)
		}
		offsets[i] = off
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(1))
		for pb.Next() {
			// pick a random pre-written offset and read
			off := offsets[r.Intn(len(offsets))]
			_, err := sm.readBlockAndDeserialize(off)
			if err != nil {
				b.Fatalf("readBlockAndDeserialize failed: %v", err)
			}
		}
	})
}
