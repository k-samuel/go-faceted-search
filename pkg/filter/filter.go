package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
)

// FilterInterface - interface for filtering realisation
type FilterInterface interface {
	GetFieldName() string
	FilterResults(facetData *index.Field, inputKeys []int64) (result []int64, err error)
}
