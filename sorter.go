package facet

import (
	"errors"
	"sort"
)

const SORT_ASC int = 0
const SORT_DESC int = 1

type SorterInterface interface {
	NewSorter(index Index) SorterInterface
	GetFieldName() string
	Sort(results []int64, field string, direction int) []int64
}

type FieldSorter struct {
	index *Index
}

func NewFieldSorter(index *Index) *FieldSorter {
	var sorter FieldSorter
	sorter.index = index
	return &sorter
}

func (sorter *FieldSorter) Sort(results []int64, field string, direction int) (result []int64, err error) {

	if !sorter.index.HasField(field) {
		err = errors.New("sort by undefined field: " + field)
		return nil, err
	}

	fieldData, _ := sorter.index.GetFieldData(field)
	s := make([]string, 0, len(fieldData.values))
	for name := range fieldData.values {
		s = append(s, name)
	}
	sort.Strings(s)

	switch direction {
	case SORT_ASC:
		sort.Sort(sort.StringSlice(s))
	case SORT_DESC:
		fallthrough
	default:
		sort.Sort(sort.Reverse(sort.StringSlice(s)))
	}
	// flip results to map
	resultsMap := make(map[int64]struct{}, len(results))
	for _, v := range results {
		resultsMap[v] = struct{}{}
	}

	resultMap := make(map[int64]struct{})

	for _, v := range s {
		if _, ok := fieldData.values[v]; ok {
			ids := intersectInt64MapKeys(fieldData.values[v].ids, resultsMap)
			if len(ids) == 0 {
				continue
			}
			for k := range ids {
				if _, ok := resultMap[k]; !ok {
					resultMap[k] = struct{}{}
				}
			}
		}
	}

	result = make([]int64, 0, len(resultMap))
	for k := range resultMap {
		result = append(result, k)
	}
	return result, err
}
