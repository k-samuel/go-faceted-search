package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
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
func (filter *ValueFilter) FilterResults(field *index.Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error) {

	var list *index.Value
	var mapLen = len(inputKeys)
	var hasInput = true
	if mapLen == 0 {
		hasInput = false
		mapLen = 100
	}

	result = make(map[int64]struct{}, mapLen)

	// collect list of record id for different values of one field
	for _, val := range filter.Values {

		if !field.HasValue(val) {
			continue
		}

		list = field.GetValue(val)
		if len(list.Ids) == 0 {
			continue
		}

		for key := range list.Ids {
			if hasInput {
				if _, ok := inputKeys[key]; ok {
					result[key] = struct{}{}
				}
			} else {
				result[key] = struct{}{}
			}
		}
	}
	return result, err
}
