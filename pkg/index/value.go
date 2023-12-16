package index

// Value - list of record id for value
type Value struct {
	Ids []int64
}

// NewValue - create value
func NewValue() *Value {
	return &Value{Ids: make([]int64, 0, 100)}
}

// addId - add record id into value struct
func (value *Value) addId(id int64) {
	value.Ids = append(value.Ids, id)
}
