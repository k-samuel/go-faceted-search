package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
)

// FilterInterface - interface for filtering realisation
type FilterInterface interface {
	GetFieldName() string
	FilterResults(facetData *index.Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error)
}
