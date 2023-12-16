package search

import (
	"math"
	"sort"

	"github.com/k-samuel/go-faceted-search/pkg/filter"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
)

// Search - faceted search
type Search struct {
	index *index.Index
}

// NewSearch Create new search instance
func NewSearch(index *index.Index) *Search {
	var search Search
	search.index = index
	return &search
}

type filterCountInfo struct {
	field string
	data  map[string]int
}

// GetIndex get index storage
func (search *Search) GetIndex() *index.Index {
	return search.index
}

// Find records using filters, limit search using list of recordId (optional)
func (search *Search) Query(q Query) (result []int64) {

	input := make(map[int64]struct{})

	if len(q.Records) > 0 {
		for _, v := range q.Records {
			input[v] = struct{}{}
		}
	}

	// Aggregates optimization for value filters.
	// The fewer elements after the first filtering, the fewer data copies and memory allocations in iterations
	if len(q.Records) == 0 && len(q.Filters) > 1 {
		q.Filters = search.sortFilters(q.Filters)
	}

	mapResult := search.findRecordsMap(q.Filters, input)

	if len(mapResult) == 0 {
		return result
	}

	result = make([]int64, 0, len(mapResult))
	for k, _ := range mapResult {
		result = append(result, k)
	}

	return result
}

/*
func (search *Search) findRecords(filters []filter.FilterInterface, inputRecords []int64) (result []int64, err error) {

	iLen := len(inputRecords)

	// return all records for empty filters
	if len(filters) == 0 {
		total := search.index.GetIdList()

		if iLen > 0 {
			return utils.IntersectSortedInt(total, inputRecords), err
		}
		result = total
		return result, err
	}

	// start value is inputRecords list
	result = inputRecords

	for _, fl := range filters {
		fieldName := fl.GetFieldName()
		if !search.index.HasField(fieldName) {
			continue
		}
		field := search.index.GetField(fieldName)
		if !field.HasValues() {
			return []int64{}, err
		}
		result, err = fl.FilterResults(field, result)
	}

	return result, err
}
*/

func (search *Search) findRecordsMap(filters []filter.FilterInterface, inputRecords map[int64]struct{}) map[int64]struct{} {

	// return all records for empty filters
	if len(filters) == 0 {
		return search.findInput(filters, inputRecords)
	}

	result := make(map[int64]struct{})

	for _, fl := range filters {
		fieldName := fl.GetFieldName()
		if !search.index.HasField(fieldName) {
			continue
		}
		field := search.index.GetField(fieldName)
		if !field.HasValues() {
			return make(map[int64]struct{})
		}

		result = fl.FilterInput(field, result)
		if len(result) == 0 {
			return result
		}
	}

	return result
}

func (search *Search) findInput(filters []filter.FilterInterface, inputRecords map[int64]struct{}) map[int64]struct{} {

	total := search.index.GetIdMap()
	if len(inputRecords) == 0 {
		return total
	}

	for k, _ := range total {
		if _, ok := inputRecords[k]; !ok {
			delete(total, k)
		}
	}

	return total
}

// AggregateFilters - find acceptable filter values
func (search *Search) Aggregate(q Aggregation) (result map[string]map[string]int) {

	result = make(map[string]map[string]int)

	input := make(map[int64]struct{})
	if len(q.Records) > 0 {
		for _, v := range q.Records {
			input[v] = struct{}{}
		}
	}

	if len(q.Filters) == 0 {
		return search.getValuesCount()
	}

	filteredRecords := make(map[int64]struct{})
	resultCache := make(map[string]map[int64]struct{})

	args := make([]filter.FilterInterface, 0, 1)

	if len(q.Filters) > 0 {
		// Aggregates optimization for value filters.
		// The fewer elements after the first filtering, the fewer data copies and memory allocations in iterations
		if len(q.Filters) > 1 {
			q.Filters = search.sortFilters(q.Filters)
		}

		for _, filter := range q.Filters {
			name := filter.GetFieldName()
			args = args[:0]
			args = append(args, filter)
			resultCache[name] = search.findRecordsMap(args, input)
		}
		filteredRecords = search.mergeFilters(resultCache, nil)

	} else if len(input) > 0 {
		filteredRecords = search.findRecordsMap([]filter.FilterInterface{}, input)
	}

	var recordIds map[int64]struct{}

	for fieldName, field := range search.index.GetFields() {

		result[fieldName] = make(map[string]int)

		if _, ok := resultCache[fieldName]; ok {
			if len(resultCache) > 1 {
				recordIds = search.mergeFilters(resultCache, &fieldName)
			} else {
				recordIds = search.findRecordsMap([]filter.FilterInterface{}, input)
			}
		} else {
			recordIds = filteredRecords
		}

		for valueName, val := range field.Values {
			result[fieldName][valueName] = utils.IntersectInt64MapKeysLen(val.Ids, recordIds)
		}
	}

	return result
}

func (search *Search) mergeFilters(cache map[string]map[int64]struct{}, skipKey *string) (result map[int64]struct{}) {
	start := true
	result = make(map[int64]struct{})

	for key, mp := range cache {
		if skipKey != nil && key == *skipKey {
			continue
		}
		if start {
			for k := range mp {
				result[k] = struct{}{}
			}
			continue
		}
		for k := range result {
			if _, ok := mp[k]; !ok {
				delete(result, k)
			}
		}
	}
	return result
}

func (search *Search) getValuesCount() (result map[string]map[string]int) {
	result = make(map[string]map[string]int)
	for fieldName, field := range search.index.GetFields() {
		result[fieldName] = make(map[string]int)
		for valName, val := range field.Values {
			result[fieldName][valName] = len(val.Ids)
		}
	}
	return result
}

type filterCount struct {
	count  int
	filter filter.FilterInterface
}

type filterValuesCount struct {
	count int
	value string
}

func (search *Search) sortFilters(filters []filter.FilterInterface) []filter.FilterInterface {

	counts := make([]*filterCount, 0, len(filters))
	var valuesInFilter int
	var valuesCount []*filterValuesCount

	// count filter values
	for index, item := range filters {

		filterCnt := &filterCount{count: math.MaxInt, filter: filters[index]}
		valFilter, ok := item.(*filter.ValueFilter)
		counts = append(counts, filterCnt)
		if !ok {
			continue
		}

		fieldName := item.GetFieldName()

		if !search.index.HasField(fieldName) {
			filterCnt.count = 0
			continue
		}
		valuesInFilter = len(valFilter.Values)
		if valuesInFilter > 1 {
			valuesCount = make([]*filterValuesCount, 0, valuesInFilter)
		}
		for _, val := range valFilter.Values {
			cnt := search.index.GetRecordsCount(fieldName, val)
			if filterCnt.count > cnt {
				filterCnt.count = cnt
			}
			if valuesInFilter > 1 {
				valuesCount = append(valuesCount, &filterValuesCount{count: cnt, value: val})
			}
		}
		if valuesInFilter > 1 {
			sort.SliceStable(valuesCount, func(i, j int) bool {
				return valuesCount[i].count < valuesCount[j].count
			})
			valFilter.Values = make([]string, 0, len(valuesCount))
			for _, v := range valuesCount {
				valFilter.Values = append(valFilter.Values, v.value)
			}
		}
	}

	sort.SliceStable(counts, func(i, j int) bool {
		return counts[i].count < counts[j].count
	})

	result := make([]filter.FilterInterface, 0, len(filters))
	for _, v := range counts {
		result = append(result, v.filter)
	}
	return result
}
