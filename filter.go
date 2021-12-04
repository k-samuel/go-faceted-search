package facet

import "strconv"

type FilterInterface interface {
	GetFieldName() string
	FilterResults(facetData *Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error)
}

type ValueFilter struct {
	FieldName string
	Values    []string
}

func (filter *ValueFilter) GetFieldName() string {
	return filter.FieldName
}

func (filter *ValueFilter) FilterResults(field *Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error) {

	var list *Value
	result = make(map[int64]struct{})

	hasInput := len(inputKeys) > 0

	// collect list for record id for different values of one field
	for _, val := range filter.Values {

		if !field.HasValue(val) {
			continue
		}

		list = field.GetValue(val)
		if len(list.ids) == 0 {
			continue
		}

		for key := range list.ids {
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

func intersectInt64MapKeys(a, b map[int64]struct{}) map[int64]struct{} {
	result := make(map[int64]struct{})
	if len(a) < len(b) {
		for key, v := range a {
			if _, ok := b[key]; ok {
				result[key] = v
			}
		}
	} else {
		for key, v := range b {
			if _, ok := a[key]; ok {
				result[key] = v
			}
		}
	}
	return result
}

func intersectInt64MapKeysLen(a, b map[int64]struct{}) int {
	var intersectLen = 0

	if len(a) < len(b) {
		for key := range a {
			if _, ok := b[key]; ok {
				intersectLen++
			}
		}
	} else {
		for key := range b {
			if _, ok := a[key]; ok {
				intersectLen++
			}
		}
	}
	return intersectLen
}

const RANGE_BOTH = 0
const RANGE_MIN = 1
const RANGE_MAX = 2

type Range struct {
	Min  float64
	Max  float64
	Type int
}

type RangeFilter struct {
	FieldName string
	Values    Range
}

func (filter *RangeFilter) GetFieldName() string {
	return filter.FieldName
}

func (filter *RangeFilter) FilterResults(field *Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error) {
	limitData := make(map[int64]struct{})
	var value float64
	// collect list for different values of one property
	for val, valObject := range field.values {
		value, err = strconv.ParseFloat(val, 8)
		if err != nil {
			return result, err
		}
		if (filter.Values.Type == RANGE_BOTH || filter.Values.Type == RANGE_MIN) && value < filter.Values.Min {
			continue
		}

		if (filter.Values.Type == RANGE_BOTH || filter.Values.Type == RANGE_MAX) && value > filter.Values.Max {
			continue
		}

		for k, v := range valObject.ids {
			limitData[k] = v
		}
	}

	if len(limitData) == 0 {
		return result, err
	}

	if len(inputKeys) == 0 {
		return limitData, err
	}

	result = intersectInt64MapKeys(limitData, inputKeys)

	return result, err
}
