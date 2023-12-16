package index


// Field - struct to store value list for index field
type Field struct {
	Values map[string]*Value
}

// NewField - create field
func NewField() *Field {
	return &Field{Values: make(map[string]*Value, 100), /*mu: &sync.Mutex{}*/}
}

// HasValues - check if field has any value
func (field *Field) HasValues() bool {
	if len(field.Values) > 0 {
		return true
	}
	return false
}

// HasValue - check if field value exists
func (field *Field) HasValue(name string) bool {
	_, ok := field.Values[name]
	return ok
}

func (field *Field) createValue(name string) *Value {
	field.Values[name] = NewValue()
	return field.Values[name]
}

// GetValue get field value by value string identifier
func (field *Field) GetValue(name string) *Value {
	return field.Values[name]
}
