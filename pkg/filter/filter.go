package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
)

// FilterInterface - interface for filtering realisation
type FilterInterface interface {
	GetFieldName() string
	FilterInput(facetData *index.Field, inputKeys map[int64]struct{}) map[int64]struct{}
}
