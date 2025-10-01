package datastruct

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"sync"
)

type Hashmap struct {
	dict  *hashmap.Map
	mutex sync.RWMutex
}

func NewHashmap() *Hashmap {
	return &Hashmap{dict: hashmap.New()}
}

func (h *Hashmap) Put(key string, value string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.dict.Put(key, value)
}

func (h *Hashmap) Get(key string) (string, bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	value, found := h.dict.Get(key)
	if found {
		return value.(string), found
	}
	return "", found
}

func (h *Hashmap) Contains(field string) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	_, found := h.dict.Get(field)
	return found
}
