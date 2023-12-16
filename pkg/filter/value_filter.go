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
func (filter *ValueFilter) FilterInput(field *index.Field, inputKeys map[int64]struct{}) map[int64]struct{} {

	if len(inputKeys) == 0 {
		return filter.FilterData(field)
	}

	var list *index.Value

	actual := make(map[int64]struct{})

	// collect list of record id for different values of one field
	for _, val := range filter.Values {

		if !field.HasValue(val) {
			continue
		}

		list = field.GetValue(val)
		if len(list.Ids) == 0 {
			continue
		}

		for _, v := range list.Ids {
			if _, ok := inputKeys[v]; ok {
				actual[v] = struct{}{}
			}
		}
	}

	return actual
}

func (filter *ValueFilter) FilterData(field *index.Field) map[int64]struct{} {
	result := make(map[int64]struct{})

	// collect list for different values of one property
	for _, val := range filter.Values {

		if !field.HasValue(val) {
			continue
		}

		list := field.GetValue(val)
		idLen := len(list.Ids)
		if idLen == 0 {
			continue
		}

		for _, v := range list.Ids {
			result[v] = struct{}{}
		}
	}
	return result
}
