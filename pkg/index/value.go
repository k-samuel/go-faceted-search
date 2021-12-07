package index

import "sync"

// Value - list of record id for value
type Value struct {
	Ids map[int64]struct{}
	mu  sync.Mutex
}

// addId - add record id into value struct
func (value *Value) addId(id int64) {
	value.mu.Lock()
	value.Ids[id] = struct{}{}
	value.mu.Unlock()
}
