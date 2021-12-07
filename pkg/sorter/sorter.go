package sorter

const SORT_ASC int = 0
const SORT_DESC int = 1

type SorterInterface interface {
	GetFieldName() string
	Sort(results []int64, field string, direction int) []int64
}
