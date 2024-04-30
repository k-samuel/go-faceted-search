package index

// Value - list of record id for value
type Value struct {
	Ids []int64
	Name string
}

// NewValue - create value
func NewValue(name string) *Value {
	return &Value{Ids: make([]int64, 0, 100), Name: name}
}

// addId - add record id into value struct
func (value *Value) addId(id int64) {
	value.Ids = append(value.Ids, id)
}
