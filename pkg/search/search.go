package search

import (
	"context"
	"github.com/k-samuel/go-faceted-search/pkg/filter"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"math"
	"runtime"
	"sort"
	"sync"
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
func (search *Search) Find(filters []filter.FilterInterface, inputRecords []int64) (result []int64, err error) {

	if len(inputRecords) > 0 {
		sort.Slice(inputRecords, func(i, j int) bool { return inputRecords[i] < inputRecords[j] })
	}

	// Aggregates optimisation for value filters.
	// The fewer elements after the first filtering, the fewer data copies and memory allocations in iterations
	if len(inputRecords) == 0 && len(filters) > 1 {
		filters = search.sortFilters(filters)
	}

	mapResult, err := search.findRecords(filters, inputRecords)
	if err != nil {
		return []int64{}, err
	}
	return mapResult, err
}

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

// AggregateFilters - find acceptable filter values
func (search *Search) AggregateFilters(filters []filter.FilterInterface, inputRecords []int64) (result map[string]map[string]int, err error) {

	if len(inputRecords) > 0 {
		sort.Slice(inputRecords, func(i, j int) bool { return inputRecords[i] < inputRecords[j] })
	}

	// Aggregates optimisation for value filters.
	// The fewer elements after the first filtering, the fewer data copies and memory allocations in iterations
	if len(inputRecords) == 0 && len(filters) > 1 {
		filters = search.sortFilters(filters)
	}

	indexedFilters := make(map[string]filter.FilterInterface, len(filters))
	indexedFilteredRecords := make([]int64, 0, 100)
	searchFields := search.index.GetFields()
	result = make(map[string]map[string]int, len(searchFields))

	if len(filters) > 0 {
		// index filters by field
		for _, filter := range filters {
			indexedFilters[filter.GetFieldName()] = filter
		}
		indexedFilteredRecords, err = search.findRecords(filters, inputRecords)
		if err != nil {
			return result, err
		}
	} else {
		if len(inputRecords) > 0 {
			indexedFilteredRecords, err = search.findRecords(filters, inputRecords)
			if err != nil {
				return result, err
			}
		}
	}

	// Create a cancel context for stopping on error
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	in := make(chan string, len(searchFields))
	out := make(chan *filterCountInfo, 10)
	errChan := make(chan error)

	go func() {
		// aggregate fields in goroutines
		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go search.aggregateField(ctx, in, out, errChan, wg, indexedFilters, indexedFilteredRecords, inputRecords)
		}
		wg.Wait()
		close(out)
	}()

	// send fields into aggregation queue
	for name := range searchFields {
		in <- name
	}
	close(in)

	// collect aggregation results
Loop:
	for {
		select {
		case err = <-errChan:
			// send cancel to goroutines on field aggregation error
			// no need to process full result
			cancel()
			// wait for goroutines stopped
			wg.Wait()
			result = make(map[string]map[string]int)
			break Loop
		case res, ok := <-out:
			if !ok {
				cancel()
				break Loop
			}
			result[res.field] = res.data
		}
	}
	return result, err
}

// aggregateField - aggregation goroutine
func (search *Search) aggregateField(
	ctx context.Context, // cancel context
	in chan string, // input channel
	out chan *filterCountInfo, // results channel
	errChan chan error, // channel for error messages
	wg *sync.WaitGroup,
	indexedFilters map[string]filter.FilterInterface, // filters indexed by field name
	indexedFilteredRecords []int64, // Total list of record id suitable for filters conditions
	inputRecords []int64, // input record id to search in
) {
	defer wg.Done()
	var filtersCopy map[string]filter.FilterInterface
	var recordIds []int64
	var field *index.Field
	var err error

	fields := search.index.GetFields()

	for {
		select {
		// cancel command
		case <-ctx.Done():
			return

		case fieldName, ok := <-in:
			if !ok {
				return
			}

			result := &filterCountInfo{field: fieldName, data: make(map[string]int)}

			field = fields[fieldName]
			if len(indexedFilters) == 0 && len(inputRecords) == 0 {
				// count values
				for val, valueObj := range field.Values {
					result.data[val] = len(valueObj.Ids)
				}
				out <- result
				runtime.Gosched()
				continue
			}

			// copy hash map
			filtersCopy = copyFilterMap(indexedFilters)

			// do not apply self filtering
			if _, ok := filtersCopy[fieldName]; ok {
				delete(filtersCopy, fieldName)
				recordIds, err = search.findRecords(extractFilters(filtersCopy), inputRecords)
				if err != nil {
					// send error (will stop other goroutines)
					errChan <- err
					return
				}
			} else {
				recordIds = indexedFilteredRecords
			}

			for vName, vList := range field.Values {
				// get records count for filter field value
				intersect := utils.IntersectCountSortedInt(vList.Ids, recordIds)
				if intersect > 0 {
					result.data[vName] = intersect
				}
			}
			out <- result
			runtime.Gosched()
		}
	}
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

func extractFilters(filters map[string]filter.FilterInterface) []filter.FilterInterface {
	var result = make([]filter.FilterInterface, 0, len(filters))
	for _, filter := range filters {
		result = append(result, filter)
	}
	return result
}

func copyFilterMap(input map[string]filter.FilterInterface) map[string]filter.FilterInterface {
	result := make(map[string]filter.FilterInterface)
	for k, v := range input {
		result[k] = v
	}
	return result
}
