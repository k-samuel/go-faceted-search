package facet

type FilterInterface interface {
	GetFieldName() string
	FilterResults(facetData *Field, inputKeys map[int64]struct{}) map[int64]struct{}
}

type ValueFilter struct {
	FieldName string
	Values    []string
}

func (filter *ValueFilter) GetFieldName() string {
	return filter.FieldName
}

func (filter *ValueFilter) FilterResults(field *Field, inputKeys map[int64]struct{}) map[int64]struct{} {

	result := make(map[int64]struct{})

	// collect list for record id for different values of one field
	for _, val := range filter.Values {
		if field.HasValue(val) {
			list := field.GetValue(val)
			if len(result) == 0 {
				result = list.ids
			} else {
				for i, _ := range list.ids {
					result[i] = struct{}{}
				}
			}
		}
	}
	// not found any or no input filtering
	if len(result) == 0 || len(inputKeys) == 0 {
		return result
	}

	// find intersect of start records and faceted results
	return intersectInt64MapKeys(result, inputKeys)
}

func intersectInt64MapKeys(a, b map[int64]struct{}) map[int64]struct{} {

	if len(a) < len(b) {
		for key, _ := range a {
			if _, ok := b[key]; !ok {
				delete(a, key)
			}
		}
		return a
	} else {
		for key, _ := range b {
			if _, ok := a[key]; !ok {
				delete(b, key)
			}
		}
		return b
	}
}
