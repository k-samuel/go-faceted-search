package sorter

import (
	"errors"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"sort"
)

// StringSorter - sorter for sorting facet data by field
type StringSorter struct {
	index *index.Index
}

// NewStringSorter - sorter constructor
func NewStringSorter(index *index.Index) *StringSorter {
	var sorter StringSorter
	sorter.index = index
	return &sorter
}

// Sort - sort faceted search results by field using index data
func (sorter *StringSorter) Sort(results []int64, field string, direction int) (result []int64, err error) {

	if !sorter.index.HasField(field) {
		err = errors.New("sort by undefined field: " + field)
		return nil, err
	}

	fieldData := sorter.index.GetField(field)
	s := make([]string, 0, len(fieldData.Values))
	for name := range fieldData.Values {
		s = append(s, name)
	}

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

	res := make(map[int64]struct{}, len(results))
	result = make([]int64, 0, len(results))

	for _, v := range s {
		if _, ok := fieldData.Values[v]; ok {
			ids := utils.IntersectInt64MapKeys(fieldData.Values[v].Ids, resultsMap)
			if len(ids) == 0 {
				continue
			}
			for k := range ids {
				if _, ok := res[k]; !ok {
					res[k] = struct{}{}
					result = append(result, k)
				}
			}
		}
	}
	return result, err
}
