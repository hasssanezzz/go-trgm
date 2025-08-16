package trgm

import (
	"testing"
)

func init() {
	Threshold = 4
}

func TestFetchFromMemoryAndDisk(t *testing.T) {
	// indexer, err := NewTrigramIndexer(t.TempDir())
	indexer, err := NewTrigramIndexer("./temp")
	if err != nil {
		t.Fatalf("Failed to create TrigramIndexer: %v", err)
	}
	defer func() {
		if err := indexer.Close(); err != nil {
			t.Fatalf("Failed to close the indexer: %v", err)
		}
	}()

	testTrigramStr := "man" // This trigram should appear in the messages below

	// Create distinct IndexEntry structs
	entries := []IndexEntry{
		{batchId: 100, offset: 1000},
		{batchId: 101, offset: 1001},
		{batchId: 102, offset: 1002},
		{batchId: 103, offset: 1003}, // 4th entry - should trigger disk write
		{batchId: 104, offset: 1004}, // Should stay in memory
		{batchId: 105, offset: 1005}, // Should stay in memory
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
