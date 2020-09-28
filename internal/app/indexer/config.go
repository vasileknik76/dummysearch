package indexer

import (
	"time"

	"github.com/vasileknik76/dummysearch/internal/app/text"
)

type IndexConfig struct {
	Language     text.Language
	UpdatePeriod time.Duration
	AutoUpdate   bool
}
