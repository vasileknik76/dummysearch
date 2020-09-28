package indexer

import (
	"time"

	"github.com/vasileknik76/dummysearch/internal/app/text"
)

type IndexConfig struct {
	Language     text.Language
	CustomIDs    bool
	UpdatePeriod time.Duration
	AutoUpdate   bool
}
