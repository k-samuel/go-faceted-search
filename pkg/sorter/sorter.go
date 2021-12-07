package sorter

// SORT_ASC - sorting order ASC
const SORT_ASC int = 0

// SORT_DESC - sorting order DESC
const SORT_DESC int = 1

// SorterInterface - interface for facet data sorters realisation
type SorterInterface interface {
	Sort(results []int64, field string, direction int) []int64
}
