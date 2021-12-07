package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"strconv"
)

// RANGE_BOTH - range type with min and max value
const RANGE_BOTH = 0

// RANGE_MIN - range type with only min value
const RANGE_MIN = 1

// RANGE_MAX - range type with only max value
const RANGE_MAX = 2

// Range value struct for RangeFilter
type Range struct {
	Min  float64
	Max  float64
	Type int
}

// RangeFilter filter facet data by field value range (numeric values)
type RangeFilter struct {
	FieldName string
	Values    Range
}

// GetFieldName - get field name
func (filter *RangeFilter) GetFieldName() string {
	return filter.FieldName
}

// FilterResults - filter facet field data
func (filter *RangeFilter) FilterResults(field *index.Field, inputKeys map[int64]struct{}) (result map[int64]struct{}, err error) {
	limitData := make(map[int64]struct{})
	var value float64
	// collect list for different values of one property
	for val, valObject := range field.Values {
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

		for k, v := range valObject.Ids {
			limitData[k] = v
		}
	}

	if len(limitData) == 0 {
		return result, err
	}

	if len(inputKeys) == 0 {
		return limitData, err
	}

	result = utils.IntersectInt64MapKeys(limitData, inputKeys)

	return result, err
}
