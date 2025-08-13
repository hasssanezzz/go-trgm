package main

import (
	"bytes"
	"os"
)

func main() {

	buf := bytes.NewBuffer(nil)
	for i := range 94 * 94 * 94 {
		buf.WriteString(intToTrigram(uint32(i)) + "\n")
	}
	os.WriteFile("result.txt", buf.Bytes(), 0666)

	// indexer := NewTrigramIndexer()
	// indexer.Index("IndexManager: An error has happended at position 45")
	// indexer.Index("IndexManager: Failed to open database connection")
	// indexer.Index("IndexManager: Failed to index batch")
	// indexer.Index("BatchManager: Failed to created batch")
	// indexer.Index("BatchManager: Failed to open batch")

	// indexer.Display()

	// results, _ := indexer.Fetch("batch")
	// for _, res := range results {
	// 	println(res)
	// }
}
