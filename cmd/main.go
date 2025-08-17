package main

import (
	"fmt"

	trgm "github.com/hasssanezzz/trigram-index"
)

func main() {
	indexer, err := trgm.NewTrigramIndexer("./temp")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := indexer.Close(); err != nil {
			panic(err)
		}
	}()

	entries, err := indexer.Fetch("wor")
	if err != nil {
		panic(err)
	}

	fmt.Println(entries)

	// if err := indexer.Index("Hello world", trgm.NewIndexEntry(5, 16)); err != nil {
	// 	panic(err)
	// }
}
