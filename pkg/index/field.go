package index

import "sync"

// Field - struct to store value list for index field
type Field struct {
	Values map[string]*Value
	mu     sync.Mutex
}

// HasValues - check if field has any value
func (field *Field) HasValues() bool {
	if len(field.Values) > 0 {
		return true
	}
	return false
}

func (field *Field) HasValue(name string) bool {
	_, ok := field.Values[name]
	return ok
}

func (field *Field) createValue(name string) *Value {
	field.mu.Lock()
	field.Values[name] = &Value{Ids: make(map[int64]struct{})}
	field.mu.Unlock()
	return field.Values[name]
}

// GetValue get field value by value string identifier
func (field *Field) GetValue(name string) *Value {
	return field.Values[name]
}
