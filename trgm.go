package gotrgm

import (
	"github.com/hasssanezzz/go-trgm/pkg/common"
	"github.com/hasssanezzz/go-trgm/pkg/core"
)

var (
	NewIndexer = core.NewTrigramIndexer
)

type (
	Indexer    = core.Indexer
	DocumentID = common.DocumentID
)
