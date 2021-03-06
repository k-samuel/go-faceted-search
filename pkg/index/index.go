package index

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

/*
 *  Data Structure
 *   Index {
 *   	ids : []  -  record id list
 *      fields: [
 *           'field1' => [] Field{
 *               values: => [
 *               	'val1' => Value{
 *                        ids: => [1,2,3,4]
 *                   },
 *                   'val2' => ...
 *               ]
 *	         },
 *           'field2 => ...
 *      ]
 *   }
 */

// Index - top level structure for facet data
type Index struct {
	fields map[string]*Field
	mu     sync.Mutex
}

// NewIndex  - Index constructor
func NewIndex() *Index {
	var index Index
	index.fields = make(map[string]*Field)
	return &index
}

// GetIdList get all record id stored in index
func (index *Index) GetIdList() []int64 {
	data := make(map[int64]struct{}, 100)
	result := make([]int64, 0, 100)
	for _, f := range index.fields {
		for _, v := range f.Values {
			for _, id := range v.Ids {
				if _, ok := data[id]; !ok {
					data[id] = struct{}{}
					result = append(result, id)
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

// GetFields get fields map
func (index *Index) GetFields() map[string]*Field {
	return index.fields
}

// Add - add record to index
func (index *Index) Add(id int64, record map[string]interface{}) {
	for key, val := range record {
		index.addValue(id, key, val)
	}
}

// HasField - check if field exists
func (index *Index) HasField(name string) bool {
	_, ok := index.fields[name]
	return ok
}

// GetRecordsCount - Get count of records registered for field value
func (index *Index) GetRecordsCount(name, value string) int {

	if _, ok := index.fields[name]; !ok {
		return 0
	}
	fld := index.fields[name]
	if _, ok := fld.Values[value]; !ok {
		return 0
	}
	return len(fld.Values[value].Ids)
}

func (index *Index) createField(name string) *Field {
	index.mu.Lock()
	index.fields[name] = NewField()
	index.mu.Unlock()
	return index.fields[name]
}

// GetField - get field struct from index
func (index *Index) GetField(name string) *Field {
	return index.fields[name]
}

// CommitChanges - save index changes
func (index *Index) CommitChanges() {
	for _, f := range index.fields {
		for _, v := range f.Values {
			sort.Slice(v.Ids, func(i, j int) bool { return v.Ids[i] < v.Ids[j] })
		}
	}
}

func (index *Index) addValue(id int64, key string, val interface{}) {
	var field *Field

	if !index.HasField(key) {
		field = index.createField(key)
	} else {
		field = index.GetField(key)
	}

	var valString string
	var value *Value

	// map
	if s, ok := val.(map[string]interface{}); ok {
		for _, v := range s {
			valString = getValueString(v)
			if !field.HasValue(valString) {
				value = field.createValue(valString)
			} else {
				value = field.GetValue(valString)
			}
			value.addId(id)
		}
		return
	}
	// array
	if s, ok := val.([]interface{}); ok {
		for _, v := range s {
			valString = getValueString(v)
			if !field.HasValue(valString) {
				value = field.createValue(valString)
			} else {
				value = field.GetValue(valString)
			}
			value.addId(id)
		}
		return
	}
	/// string
	valString = getValueString(val)
	if !field.HasValue(valString) {
		value = field.createValue(valString)
	} else {
		value = field.GetValue(valString)
	}
	value.addId(id)
}

// getValueString - convert value to string
func getValueString(val interface{}) string {

	if s, ok := val.(bool); ok {
		if s {
			return "1"
		}
		return "0"
	}

	if s, ok := val.(string); ok {
		return s
	}

	if s, ok := val.(int64); ok {
		return strconv.FormatInt(s, 10)
	}

	if s, ok := val.(int); ok {
		return strconv.Itoa(s)
	}

	if s, ok := val.(float64); ok {
		return strconv.FormatFloat(s, 'f', -1, 64)
	}

	fmt.Printf("undefined value type %T->%q\n", val, val)
	panic("undefined value type")
}
