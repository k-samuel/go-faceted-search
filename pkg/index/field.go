package index

import (
	"sort"
)

// Field - struct to store value list for index field
type Field struct {
	Values   []*Value
	valueMap map[string]*Value
}

// NewField - create field
func NewField() *Field {
	return &Field{Values: make([]*Value, 0, 10), valueMap: make(map[string]*Value, 10)}
}

// HasValues - check if field has any value
func (field *Field) HasValues() bool {
	return len(field.Values) > 0
}

func (field *Field) createValue(name string) *Value {
	val := NewValue(name)
	field.Values = append(field.Values, val)
	field.valueMap[name] = val
	return val
}

// GetValue get field value by value string identifier
func (field *Field) GetValue(name string) (val *Value, ok bool) {
	val, ok = field.valueMap[name]
	return val, ok
}

func (field *Field) SortValues() {
	for _, v := range field.Values {
		sort.Slice(v.Ids, func(i, j int) bool { return v.Ids[i] < v.Ids[j] })
	}
	sort.Slice(field.Values, func(i, j int) bool { return len(field.Values[i].Ids) < len(field.Values[j].Ids) })
}
