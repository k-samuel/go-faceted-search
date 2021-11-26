package facet

type Search struct {
	index *Index
}

// NewSearch Create new search instance
func NewSearch(index *Index) *Search {
	var search Search
	search.index = index
	return &search
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
	var recordIds map[int64]struct{}

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
	var filtersCopy map[string]FilterInterface

	for name, field := range search.index.fields {

		if len(indexedFilters) == 0 && len(inputRecords) == 0 {
			// count values
			for val, valueObj := range field.values {
				if _, ok := result[name]; !ok {
					result[name] = make(map[string]int)
				}
				result[name][val] = len(valueObj.ids)
			}
			continue
		}

		// copy hash map
		filtersCopy = copyFilterMap(indexedFilters)

		// do not apply self filtering
		if _, ok := filtersCopy[name]; ok {
			delete(filtersCopy, name)
			found, err := search.Find(extractFilters(filtersCopy), inputRecords)
			if err != nil {
				return result, err
			}
			recordIds = flipInt64ToMap(found)
		} else {
			recordIds = indexedFilteredRecords
		}

		for vName, vList := range field.values {
			// need to count values
			intersect := intersectInt64MapKeys(vList.ids, recordIds)
			if len(intersect) > 0 {
				if _, ok := result[name]; !ok {
					result[name] = make(map[string]int)
				}
				result[name][vName] = len(intersect)
			}
		}
	}
	return result, err
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
