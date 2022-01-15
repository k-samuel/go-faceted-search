package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
)

// ValueFilter - filter facet data by field value
type ValueFilter struct {
	FieldName string
	Values    []string
}

// GetFieldName - get field name
func (filter *ValueFilter) GetFieldName() string {
	return filter.FieldName
}

// FilterResults - filter facet field data
func (filter *ValueFilter) FilterResults(field *index.Field, inputKeys []int64) (result []int64, err error) {

	var list *index.Value
	var mapLen = len(inputKeys)
	var hasInput = true
	if mapLen == 0 {
		hasInput = false
		mapLen = 100
	}

	result = make([]int64, 0, mapLen)

	// collect list of record id for different values of one field
	for _, val := range filter.Values {

		if !field.HasValue(val) {
			continue
		}

		list = field.GetValue(val)
		if len(list.Ids) == 0 {
			continue
		}

		if hasInput {
			result = append(result, utils.IntersectSortedInt(list.Ids, inputKeys)...)
		} else {
			result = append(result, list.Ids...)
		}
	}
	if len(result) > 1 {
		result = utils.Deduplicate(result)
	}
	return result, err
}
