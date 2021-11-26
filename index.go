package facet

import (
	"fmt"
	"strconv"
	"sync"
)

type Index struct {
	fields map[string]*Field
	ids    []int64
	mu     sync.Mutex
}

func NewIndex() *Index {
	var index Index
	index.fields = make(map[string]*Field)
	index.ids = make([]int64, 0, 500)
	return &index
}

func (index *Index) GetData() map[string]*Field {
	return index.fields
}

func (index *Index) Add(id int64, record map[string]interface{}) {
	index.mu.Lock()
	index.ids = append(index.ids, id)
	index.mu.Unlock()
	for key, val := range record {
		index.addValue(id, key, val)
	}
}

func (index *Index) HasField(name string) bool {
	_, ok := index.fields[name]
	return ok
}

func (index *Index) createField(name string) *Field {
	index.mu.Lock()
	index.fields[name] = &Field{values: make(map[string]*Value)}
	index.mu.Unlock()
	return index.fields[name]
}

func (index *Index) GetField(name string) *Field {
	return index.fields[name]
}

func (index *Index) GetItemsCount() int {
	return len(index.ids)
}

func (index *Index) GetFieldData(fieldName string) (val *Field, ok bool) {
	if _, ok := index.fields[fieldName]; ok {
		return index.fields[fieldName], true
	}
	return val, false
}

// GetAllRecordId get all record id stored in index
func (index *Index) GetAllRecordId() []int64 {
	return index.ids
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

func getValueString(val interface{}) string {

	if s, ok := val.(bool); ok {
		if s {
			return "1"
		} else {
			return "0"
		}
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

type Field struct {
	values map[string]*Value
	mu     sync.Mutex
}

func (field *Field) HasValues() bool {
	if len(field.values) > 0 {
		return true
	}
	return false
}

func (field *Field) HasValue(name string) bool {
	_, ok := field.values[name]
	return ok
}

func (field *Field) createValue(name string) *Value {
	field.mu.Lock()
	field.values[name] = &Value{ids: make(map[int64]struct{})}
	field.mu.Unlock()
	return field.values[name]
}

func (field *Field) GetValue(name string) *Value {
	return field.values[name]
}

type Value struct {
	ids map[int64]struct{}
	mu  sync.Mutex
}

func (value *Value) addId(id int64) {
	value.mu.Lock()
	value.ids[id] = struct{}{}
	value.mu.Unlock()
}
