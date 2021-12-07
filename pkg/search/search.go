package search

import (
	"context"
	"github.com/k-samuel/go-faceted-search/pkg/filter"
	"github.com/k-samuel/go-faceted-search/pkg/index"
	"github.com/k-samuel/go-faceted-search/pkg/utils"
	"runtime"
	"sync"
)

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

// Find records using filters, limit search using list of recordId (optional)
func (search *Search) Find(filters []filter.FilterInterface, inputRecords []int64) (result []int64, err error) {

	input := make(map[int64]struct{})
	if len(inputRecords) > 0 {
		input = utils.FlipInt64ToMap(inputRecords)
	}

	mapResult, err := search.findRecords(filters, input)
	if err != nil {
		return []int64{}, err
	}
	// Convert result map into array of int
	resLen := len(mapResult)
	if resLen > 0 {
		result = make([]int64, 0, resLen)
		for key := range mapResult {
			result = append(result, key)
		}
	}
	return result, err
}

func (search *Search) findRecords(filters []filter.FilterInterface, inputRecords map[int64]struct{}) (result map[int64]struct{}, err error) {

	result = make(map[int64]struct{})
	iLen := len(inputRecords)

	// return all records for empty filters
	if len(filters) == 0 {
		total := search.index.GetIdList()
		if iLen > 0 {
			return utils.IntersectRecAndMapKeysToMap(total, inputRecords), err
		}
		for _, v := range total {
			result[v] = struct{}{}
		}
		return result, err
	}

	// start value is inputRecords list
	result = inputRecords
	for _, filter := range filters {
		fieldName := filter.GetFieldName()
		if !search.index.HasField(fieldName) {
			continue
		}
		field := search.index.GetField(fieldName)
		if !field.HasValues() {
			return result, err
		}
		result, err = filter.FilterResults(field, result)
	}
	return result, err
}

// AggregateFilters - find acceptable filter values
func (search *Search) AggregateFilters(filters []filter.FilterInterface, inputRecords []int64) (result map[string]map[string]int, err error) {

	input := make(map[int64]struct{})
	if len(inputRecords) > 0 {
		input = utils.FlipInt64ToMap(inputRecords)
	}

	result = make(map[string]map[string]int)
	indexedFilters := make(map[string]filter.FilterInterface)

	indexedFilteredRecords := make(map[int64]struct{})

	if len(filters) > 0 {
		// index filters by field
		for _, filter := range filters {
			indexedFilters[filter.GetFieldName()] = filter
		}
		indexedFilteredRecords, err = search.findRecords(filters, input)
		if err != nil {
			return result, err
		}
	}

	// Create a cancel context for stopping on error
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	in := make(chan string, 10)
	out := make(chan *filterCountInfo, 10)
	errChan := make(chan error)

	go func() {
		// aggregate fields in goroutines
		for i := 0; i < runtime.NumCPU(); i++ {
			wg.Add(1)
			go search.aggregateField(ctx, in, out, errChan, wg, indexedFilters, indexedFilteredRecords, input)
		}
		wg.Wait()
		close(out)
	}()

	// send fields into aggregation queue
	for name := range search.index.GetFields() {
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
	indexedFilteredRecords map[int64]struct{}, // Total list of record id suitable for filters conditions
	inputRecords map[int64]struct{}, // input record id to search in
) {
	defer wg.Done()
	var filtersCopy map[string]filter.FilterInterface
	var recordIds map[int64]struct{}
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
				intersect := utils.IntersectInt64MapKeysLen(vList.Ids, recordIds)
				if intersect > 0 {
					result.data[vName] = intersect
				}
			}
			out <- result
			runtime.Gosched()
		}
	}
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
