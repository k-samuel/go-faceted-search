package index

import "sync"

// Value - list of record id for value
type Value struct {
	Ids []int64
	mu  sync.Mutex
}

// addId - add record id into value struct
func (value *Value) addId(id int64) {
	value.mu.Lock()
	value.Ids = append(value.Ids, id)
	value.mu.Unlock()
}
