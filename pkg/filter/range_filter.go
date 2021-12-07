package filter

import (
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"strconv"
)

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
