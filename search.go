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
func (search *Search) Find(filters []FilterInterface, inputRecords []int64) []int64 {

	var idFilter map[int64]struct{}
	var result []int64

	// convert inputRecords into hash map for fast search
	iLen := len(inputRecords)
	if iLen > 0 {
		idFilter = make(map[int64]struct{}, iLen)
		for _, val := range inputRecords {
			idFilter[val] = struct{}{}
		}
	}

	// return all records for empty filters
	if len(filters) == 0 {
		total := search.index.GetAllRecordId()
		if iLen > 0 {
			return intersectRecAndMapKeys(total, idFilter)
		}
		return total
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
			return result
		}
		mapResult = filter.FilterResults(field, mapResult)
	}

	// Convert result map into array of int
	resLen := len(mapResult)
	if resLen > 0 {
		result = make([]int64, 0, resLen)
		for key := range mapResult {
			result = append(result, key)
		}
	}
	return result
}

// AggregateFilters - find acceptable filter values
func (search *Search) AggregateFilters(filters []FilterInterface, inputRecords []int64) map[string]map[string]int {

	result := make(map[string]map[string]int)
	indexedFilters := make(map[string]FilterInterface)

	var filteredRecords []int64
	var recordIds map[int64]struct{}

	indexedFilteredRecords := make(map[int64]struct{})

	if len(filters) > 0 {
		// index filters by field
		for _, filter := range filters {
			indexedFilters[filter.GetFieldName()] = filter
		}
		filteredRecords = search.Find(filters, inputRecords)
		// flip filtered records
		if len(filteredRecords) > 0 {
			for _, v := range filteredRecords {
				indexedFilteredRecords[v] = struct{}{}
			}
			//filteredRecords = array_flip($filteredRecords);
		}
	}
	for name, field := range search.index.fields {
		if len(indexedFilters) == 0 && len(inputRecords) == 0 {
			// count values
			for val, valueObj := range field.values {
				if _, ok := result[name]; !ok {
					result[name] = make(map[string]int)
				}
				result[name][val] = len(valueObj.ids)
			}
		} else {
			filtersCopy := indexedFilters
			// do not apply self filtering
			if _, ok := filtersCopy[name]; ok {
				delete(filtersCopy, name)
				recordIds = flipInt64ToMap(search.Find(extractFilters(filtersCopy), inputRecords))
			} else {
				recordIds = indexedFilteredRecords
			}

			for vName, v := range field.values {
				// need to count values
				intersect := intersectInt64MapKeys(v.ids, recordIds)
				if len(intersect) > 0 {
					if _, ok := result[name]; !ok {
						result[name] = make(map[string]int)
					}
					result[name][vName] = len(intersect)
				}
			}
		}
	}
	return result
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
