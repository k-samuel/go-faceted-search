package facet

type Search struct {
	index *Index
}

func NewSearch(index *Index) *Search {
	var search Search
	search.index = index
	return &search
}

func (search *Search) Find(filters []FilterInterface, inputRecords []int64) []int64 {

	var idFilter map[int64]struct{}
	var result []int64

	// конвертируем список inputRecords в map для быстрого поиска
	iLen := len(inputRecords)
	if iLen > 0 {
		idFilter = make(map[int64]struct{}, iLen)
		for _, val := range inputRecords {
			idFilter[val] = struct{}{}
		}
	}

	// если не переданы фильтры
	if len(filters) == 0 {
		total := search.index.GetAllRecordId()
		if iLen > 0 {
			return intersectRecAndMapKeys(total, idFilter)
		}
		return total
	}

	var mapResult map[int64]struct{}

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
		for key, _ := range mapResult {
			result = append(result, key)
		}
	}
	return result
}

func (index *Index) GetAllRecordId() []int64 {
	return index.ids
}

// Intersection of records ids and filter list
func intersectRecAndMapKeys(records []int64, keys map[int64]struct{}) []int64 {
	result := make([]int64, 0, len(keys))
	for _, v := range records {
		if _, ok := keys[v]; ok {
			result = append(result, v)
		}
	}
	return result
}
