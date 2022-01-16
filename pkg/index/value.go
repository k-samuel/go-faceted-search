package index

import "sync"

// Value - list of record id for value
type Value struct {
	mu  *sync.Mutex
	Ids []int64
}

// NewValue - create value
func NewValue() *Value {
	return &Value{Ids: make([]int64, 0, 100), mu: &sync.Mutex{}}
}

// addId - add record id into value struct
func (value *Value) addId(id int64) {
	value.mu.Lock()
	value.Ids = append(value.Ids, id)
	value.mu.Unlock()
}
