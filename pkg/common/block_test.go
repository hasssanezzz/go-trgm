package common

import (
	"bytes"
	"math/rand"
	"testing"
)

func compareEntries(a, b []DocumentID) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}

	return true
}

func TestNewBlock(t *testing.T) {
	n := 1000
	tri := uint32(69)
	entries := make([]DocumentID, n)
	for i := range n {
		entries[i] = [DocumentIdSize]byte{byte(i), byte(rand.Intn(256))}
	}

	block := NewBlock(tri, entries)

	t.Run("Before serialization", func(t *testing.T) {
		if !compareEntries(block.Entries(), entries) {
			t.Fatal("entry slice mismatch")
		}

		if tri != block.tri {
			t.Fatal("tri mismatch")
		}

		block.Delete(entries[0])
		block.DeleteByIndex(1)

		if len(block.Entries()) != n-2 {
			t.Fatal("entry length mismatch after deleting 2 entries")
		}

		if !compareEntries(block.Entries(), entries[2:]) {
			t.Fatal("entry array mismatch")
		}
	})

	t.Run("After serialization", func(t *testing.T) {
		newBlock, err := BlockFromReader(bytes.NewReader(block.Serialize()))
		if err != nil {
			t.Fatalf("failed to deserialize block: %v", err)
		}

		if newBlock.tri != tri {
			t.Fatal("tri mismatch")
		}

		if len(newBlock.Entries()) != n-2 {
			t.Fatalf("entry length mismatch after deserialization, want %d got %d", n-2, len(newBlock.Entries()))
		}

		if !compareEntries(newBlock.Entries(), entries[2:]) {
			t.Fatal("entry array mismatch after deserialization")
		}
	})

}
