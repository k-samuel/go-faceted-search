package utils

func IntersectInt64MapKeys(a, b map[int64]struct{}) map[int64]struct{} {
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

func IntersectInt64MapKeysLen(a, b map[int64]struct{}) int {
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

func FlipInt64ToMap(list []int64) map[int64]struct{} {
	result := make(map[int64]struct{})
	for _, v := range list {
		result[v] = struct{}{}
	}
	return result
}

func CopyInt64Map(input map[int64]struct{}) map[int64]struct{} {
	result := make(map[int64]struct{})
	for k, v := range input {
		result[k] = v
	}
	return result
}

// IntersectRecAndMapKeys Intersection of records ids and filter list
func IntersectRecAndMapKeys(records []int64, keys map[int64]struct{}) []int64 {
	result := make([]int64, 0, len(keys))
	for _, v := range records {
		if _, ok := keys[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// IntersectRecAndMapKeysToMap Intersection of records ids and filter list
func IntersectRecAndMapKeysToMap(records []int64, keys map[int64]struct{}) map[int64]struct{} {
	result := make(map[int64]struct{})
	for _, v := range records {
		if _, ok := keys[v]; ok {
			result[v] = struct{}{}
		}
	}
	return result
}
