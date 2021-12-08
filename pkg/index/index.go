package index

import (
	"fmt"
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
	ids    []int64
	mu     sync.Mutex
}

// NewIndex  - Index constructor
func NewIndex() *Index {
	var index Index
	index.fields = make(map[string]*Field)
	index.ids = make([]int64, 0, 500)
	return &index
}

// GetIdList get all record id stored in index
func (index *Index) GetIdList() []int64 {
	data := make([]int64, len(index.ids))
	copy(data, index.ids)
	return data
}

// Ids get pointer to list of record id (unsafe)
func (index *Index) Ids() []int64 {
	return index.ids
}

// GetFields get fields map
func (index *Index) GetFields() map[string]*Field {
	return index.fields
}

// Add - add record to index
func (index *Index) Add(id int64, record map[string]interface{}) {
	index.mu.Lock()
	index.ids = append(index.ids, id)
	index.mu.Unlock()
	for key, val := range record {
		index.addValue(id, key, val)
	}
}

// HasField - check if field exists
func (index *Index) HasField(name string) bool {
	_, ok := index.fields[name]
	return ok
}

func (index *Index) createField(name string) *Field {
	index.mu.Lock()
	index.fields[name] = &Field{Values: make(map[string]*Value)}
	index.mu.Unlock()
	return index.fields[name]
}

// GetField - get field struct from index
func (index *Index) GetField(name string) *Field {
	return index.fields[name]
}

// GetItemsCount - get total count records in index
func (index *Index) GetItemsCount() int {
	return len(index.ids)
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
