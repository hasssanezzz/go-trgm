package core

import (
	"os"
	"strings"
	"testing"

	"github.com/hasssanezzz/go-trgm/pkg/common"
)

var primes = []int{505447,
	505447,
	511111,
	524287,
	524287,
	541027,
	541027,
	563417,
	584141,
	584141,
	593933,
	593933,
	593993,
	593993,
	665557,
	667333}

func init() {
	Threshold = 4
}

func readLines() []string {
	path := "../../testing/0.5mil.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(data), "\n")
}

func makeDocumentId(n int) common.DocumentID {
	d := common.DocumentID{}
	for i := range common.DocumentIdSize {
		d[i] = uint8((n + 1) * primes[i] % 256)
	}
	return d
}

func BenchmarkIndexer(b *testing.B) {
	indexer, err := NewTrigramIndexer(b.TempDir())
	if err != nil {
		b.Fatalf("Failed to create TrigramIndexer: %v", err)
	}
	defer func() {
		if err := indexer.Close(); err != nil {
			b.Fatalf("Failed to close the indexer: %v", err)
		}
	}()

	lines := readLines()

	b.Run("Indexing", func(b *testing.B) {
		b.ResetTimer()
		for i := range b.N {
			if err := indexer.Index(lines[i%len(lines)], makeDocumentId(i)); err != nil {
				panic(err)
			}
		}
	})

	b.Run("Fetching", func(b *testing.B) {
		b.ResetTimer()
		for i := range b.N {
			if _, err := indexer.Fetch(lines[i%len(lines)]); err != nil {
				panic(err)
			}
		}
	})
}

func TestFetchFromMemoryAndDisk(t *testing.T) {
	indexer, err := NewTrigramIndexer(t.TempDir())
	if err != nil {
		t.Fatalf("Failed to create TrigramIndexer: %v", err)
	}
	defer func() {
		if err := indexer.Close(); err != nil {
			t.Fatalf("Failed to close the indexer: %v", err)
		}
	}()

	testTrigramStr := "man" // This trigram should appear in the messages below

	// Create distinct common.IndexEntry structs
	entries := []common.DocumentID{
		makeDocumentId(100),
		makeDocumentId(101),
		makeDocumentId(102),
		makeDocumentId(103), // 4th entry - should trigger disk write
		makeDocumentId(104), // Should stay in memory
		makeDocumentId(105), // Should stay in memory
	}

	messages := []string{
		"Manager: failed to create database",    // Contains "man"
		"BatchManager: Error at position 23",    // Contains "man"
		"DatabaseManager: connection could not", // Contains "man"
		"CacheManager: Failed to process file",  // Contains "man"
		"TaskManager: timeout while update",     // Contains "man" (5th)
		"QueueManager: exception during delete", // Contains "man" (6th)
	}

	for i, msg := range messages {
		err := indexer.Index(msg, entries[i])
		if err != nil {
			t.Fatalf("Failed to index message %d ('%s'): %v", i, msg, err)
		}
	}

	fetchedEntries, err := indexer.Fetch(testTrigramStr)
	if err != nil {
		t.Fatalf("Fetch failed for pattern '%s': %v", testTrigramStr, err)
	}

	if len(fetchedEntries) != 6 {
		t.Fatalf("Fetch expected count mismatch wanted %d got %d", 6, len(fetchedEntries))
	}
}
