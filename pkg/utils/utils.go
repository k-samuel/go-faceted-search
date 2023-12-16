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

// IntersectSortedInt intersect sorted int slices
func IntersectSortedInt(a, b []int64) []int64 {
	if len(a) == 0 || len(b) == 0 {
		return []int64{}
	}

	compareCount := len(b)
	comparePointer := 0

	result := make([]int64, 0, 100)

	for _, value := range a {

		if comparePointer >= compareCount {
			break
		}

		if value < b[comparePointer] {
			continue
		}
		for ; comparePointer < compareCount; comparePointer++ {
			if b[comparePointer] < value {
				continue
			}

			if b[comparePointer] == value {
				result = append(result, value)
				break
			}

			if b[comparePointer] > value {
				break
			}
		}
	}
	return result
}

func IntersectReplaceSortedInt(source, target []int64) []int64 {
	var i, size int
	var has bool
	var min int64 = -1
	size = len(target)

	for i = 0; i < size; i++ {
		has = false
		v := target[i]
		if v == min {
			break
		}
		for _, val := range source {
			if val > v {
				break
			}
			if val == v {
				has = true
				break
			}
		}
		if !has {
			if min == -1 {
				min = v
			}
			// Remove the element at index i from target.
			target[i] = target[len(target)-1] // Copy last element to index i.
			//target[len(target)-1] = ""   // Erase last element (write zero value).
			target = target[:len(target)-1] // Truncate slice.
			i--                             // retry
			size--
		}
	}
	sort.Slice(target, func(i, j int) bool { return target[i] < target[j] })
	return target
}

// IntersectCountSortedInt get intersect count for sorted int slices
func IntersectCountSortedInt(a, b []int64) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	result := 0

	compareCount := len(b)
	comparePointer := 0

	for _, value := range a {
		if comparePointer >= compareCount {
			break
		}
		if value < b[comparePointer] {
			continue
		}
		for ; comparePointer < compareCount; comparePointer++ {
			if b[comparePointer] < value {
				continue
			}

			if b[comparePointer] == value {
				result++
				break
			}

			if b[comparePointer] > value {
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
