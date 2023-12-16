package search

import (
	"github.com/k-samuel/go-faceted-search/pkg/filter"
)
// Aggregation - query builder for records aggregation
type Aggregation struct {
	Filters []filter.FilterInterface
	Records []int64
	CountRecords   bool
}