package maps

import (
	"github.com/emirpasic/gods/maps/hashmap"
)

type Hashmap struct {
	dict *hashmap.Map
}

func NewHashmap() *Hashmap {
	return &Hashmap{dict: hashmap.New()}
}

func (h *Hashmap) Put(key string, value string) {
	h.dict.Put(key, value)
}

func (h *Hashmap) Get(key string) (string, bool) {
	value, found := h.dict.Get(key)
	if found {
		return value.(string), found
	}
	return "", found
}

func (h *Hashmap) Contains(key string) bool {
	_, found := h.dict.Get(key)
	return found
}
