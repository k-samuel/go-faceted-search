package utils

import "sort"

// IntersectInt64MapKeys - intersection of int64 maps
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

// IntersectInt64MapKeysLen - get count of intersected values
func IntersectInt64MapKeysLen(a []int64, b map[int64]struct{}) int {
	var intersectLen = 0
	for _, key := range a {
		if _, ok := b[key]; ok {
			intersectLen++
		}
	}
	return intersectLen
}

// FlipInt64ToMap - convert in64 slice into map as keys
func FlipInt64ToMap(list []int64) map[int64]struct{} {
	result := make(map[int64]struct{})
	for _, v := range list {
		result[v] = struct{}{}
	}
	return result
}

// CopyInt64Map - copy map
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
	result := make(map[int64]struct{}, len(keys))
	for _, v := range records {
		if _, ok := keys[v]; ok {
			result[v] = struct{}{}
		}
	}
	return result
}

// IntersectSortedInt intersect sorted int slices
func IntersectSortedInt(a, b []int64) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return []int64{}
	}
	var start []int64
	var compare []int64

	aLen := len(a)
	bLen := len(b)
	// chose minimal slice
	if aLen < bLen {
		start = a
		compare = b
	} else {
		start = b
		compare = a
	}
	compareCount := len(compare)
	comparePointer := 0

	result := make([]int64, 0, 100)

	for _, value := range start {

		if comparePointer >= compareCount {
			break
		}

		if value < compare[comparePointer] {
			continue
		}
		for ; comparePointer < compareCount; comparePointer++ {
			if compare[comparePointer] < value {
				continue
			}

			if compare[comparePointer] == value {
				result = append(result, value)
				break
			}

			if compare[comparePointer] > value {
				break
			}
		}
	}
	return result
}

// IntersectCountSortedInt get intersect count for sorted int slices
func IntersectCountSortedInt(a, b []int64) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	result := 0
	var start []int64
	var compare []int64

	aLen := len(a)
	bLen := len(b)
	// chose minimal slice
	if aLen < bLen {
		start = a
		compare = b
	} else {
		start = b
		compare = a
	}
	compareCount := len(compare)
	comparePointer := 0

	for _, value := range start {
		if comparePointer >= compareCount {
			break
		}
		if value < compare[comparePointer] {
			continue
		}
		for ; comparePointer < compareCount; comparePointer++ {
			if compare[comparePointer] < value {
				continue
			}

			if compare[comparePointer] == value {
				result++
				break
			}

			if compare[comparePointer] > value {
				break
			}
		}
	}
	return result
}

// Deduplicate - remove duplicates from int slice
func Deduplicate(in []int64) []int64 {
	sort.Slice(in, func(i, j int) bool { return in[i] < in[j] })
	// In-place deduplicate https://github.com/golang/go/wiki/SliceTricks
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		// preserve the original data
		// in[i], in[j] = in[j], in[i]
		// only set what is required
		in[j] = in[i]
	}
	result := in[:j+1]
	return result
}
