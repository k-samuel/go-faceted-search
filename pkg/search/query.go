package search

import (
	"github.com/k-samuel/go-faceted-search/pkg/filter"
)

// Query - query builder for records search
type Query struct {
	Filters []filter.FilterInterface
	Limit   int
	Records []int64
}
