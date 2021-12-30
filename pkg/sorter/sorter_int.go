package sorter

import (
	"errors"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"sort"
	"strconv"
)

// IntSorter - sorter for sorting facet data by field
type IntSorter struct {
	index *index.Index
}

// NewIntSorter - sorter constructor
func NewIntSorter(index *index.Index) *IntSorter {
	var sorter IntSorter
	sorter.index = index
	return &sorter
}

// Sort - sort faceted search results by field using index data
func (sorter *IntSorter) Sort(results []int64, field string, direction int) (result []int64, err error) {

	if !sorter.index.HasField(field) {
		err = errors.New("sort by undefined field: " + field)
		return nil, err
	}

	var val int
	var str string

	fieldData := sorter.index.GetField(field)
	s := make([]int, 0, len(fieldData.Values))
	for name := range fieldData.Values {
		val, err = strconv.Atoi(name)
		if err != nil {
			return result, err
		}
		s = append(s, val)
	}

	switch direction {
	case SORT_ASC:
		sort.Sort(sort.IntSlice(s))
	case SORT_DESC:
		fallthrough
	default:
		sort.Sort(sort.Reverse(sort.IntSlice(s)))
	}

	// flip results to map
	resultsMap := make(map[int64]struct{}, len(results))
	for _, v := range results {
		resultsMap[v] = struct{}{}
	}

	//res := make(map[int64]struct{}, len(results))
	result = make([]int64, 0, len(results))

	for _, v := range s {
		str = strconv.Itoa(v)
		if err != nil {
			return result, err
		}
		if _, ok := fieldData.Values[str]; ok {
			ids := utils.IntersectRecAndMapKeys(fieldData.Values[str].Ids, resultsMap)
			if len(ids) == 0 {
				continue
			}
			for _, k := range ids {
				result = append(result, k)
				delete(resultsMap, k)
			}
		}
	}

	return result, err
}
