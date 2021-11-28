package facet

import (
	"context"
	"runtime"
	"sync"
)

type Search struct {
	index *Index
}

// NewSearch Create new search instance
func NewSearch(index *Index) *Search {
	var search Search
	search.index = index
	return &search
}

type filterCountInfo struct {
	field string
	data  map[string]int
}

// Find records using filters, limit search using list of recordId (optional)
func (search *Search) Find(filters []FilterInterface, inputRecords []int64) (result []int64, err error) {

	var idFilter = make(map[int64]struct{}, 0)
	result = make([]int64, 0, 10)

	// convert inputRecords into hash map for fast search
	iLen := len(inputRecords)
	if iLen > 0 {
		idFilter = flipInt64ToMap(inputRecords)
	}

	// return all records for empty filters
	if len(filters) == 0 {
		total := search.index.GetAllRecordId()
		if iLen > 0 {
			return intersectRecAndMapKeys(total, idFilter), err
		}
		return total, err
	}

	var mapResult map[int64]struct{}

	// start value is inputRecords list
	mapResult = idFilter

	for _, filter := range filters {
		fieldName := filter.GetFieldName()
		if !search.index.HasField(fieldName) {
			continue
		}
		field := search.index.GetField(fieldName)
		if !field.HasValues() {
			return result, err
		}
		mapResult, err = filter.FilterResults(field, mapResult)
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

// AggregateFilters - find acceptable filter values
func (search *Search) AggregateFilters(filters []FilterInterface, inputRecords []int64) (result map[string]map[string]int, err error) {

	result = make(map[string]map[string]int)
	indexedFilters := make(map[string]FilterInterface)

	var filteredRecords []int64

	indexedFilteredRecords := make(map[int64]struct{})

	if len(filters) > 0 {
		// index filters by field
		for _, filter := range filters {
			indexedFilters[filter.GetFieldName()] = filter
		}
		filteredRecords, err = search.Find(filters, inputRecords)
		// flip filtered records
		if len(filteredRecords) > 0 {
			indexedFilteredRecords = flipInt64ToMap(filteredRecords)
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
			go aggregateField(in, out, ctx, errChan, wg, indexedFilters, indexedFilteredRecords, inputRecords, search)
		}
		wg.Wait()
		close(out)
	}()

	// send fields into aggregation queue
	for name := range search.index.fields {
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
func aggregateField(
	in chan string, // input channel
	out chan *filterCountInfo, // results channel
	ctx context.Context, // cancel context
	errChan chan error, // channel for error messages
	wg *sync.WaitGroup,
	indexedFilters map[string]FilterInterface, // filters indexed by field name
	indexedFilteredRecords map[int64]struct{}, // Total list of record id suitable for filters conditions
	inputRecords []int64, // input record id to search in
	search *Search, // search object

) {
	defer wg.Done()
	var filtersCopy map[string]FilterInterface
	var recordIds map[int64]struct{}
	var field *Field

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

			field = search.index.fields[fieldName]
			if len(indexedFilters) == 0 && len(inputRecords) == 0 {
				// count values
				for val, valueObj := range field.values {
					result.data[val] = len(valueObj.ids)
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
				found, err := search.Find(extractFilters(filtersCopy), inputRecords)
				if err != nil {
					// send error (will stop other goroutines)
					errChan <- err
					return
				}
				recordIds = flipInt64ToMap(found)
			} else {
				recordIds = indexedFilteredRecords
			}

			for vName, vList := range field.values {
				// get records count for filter field value
				intersect := intersectInt64MapKeys(vList.ids, recordIds)
				if len(intersect) > 0 {
					result.data[vName] = len(intersect)
				}
			}
			out <- result
			runtime.Gosched()
		}
	}
}

func extractFilters(filters map[string]FilterInterface) []FilterInterface {
	var result = make([]FilterInterface, 0, len(filters))
	for _, filter := range filters {
		result = append(result, filter)
	}
	return result
}

func flipInt64ToMap(list []int64) map[int64]struct{} {
	result := make(map[int64]struct{})
	for _, v := range list {
		result[v] = struct{}{}
	}
	return result
}

func copyInt64Map(input map[int64]struct{}) map[int64]struct{} {
	result := make(map[int64]struct{})
	for k, v := range input {
		result[k] = v
	}
	return result
}

func copyFilterMap(input map[string]FilterInterface) map[string]FilterInterface {
	result := make(map[string]FilterInterface)
	for k, v := range input {
		result[k] = v
	}
	return result
}

// intersectRecAndMapKeys Intersection of records ids and filter list
func intersectRecAndMapKeys(records []int64, keys map[int64]struct{}) []int64 {
	result := make([]int64, 0, len(keys))
	for _, v := range records {
		if _, ok := keys[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
