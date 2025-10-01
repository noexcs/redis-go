package datastruct

import (
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

// Sortedset https://redis.io/docs/data-types/sorted-sets/
type Sortedset struct {
	sortedset *treemap.Map
	mutex     sync.Mutex
}

func (s *Sortedset) Add(score int64, elem string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sortedset.Put(score, elem)
}

func NewSortedset() *Sortedset {
	return &Sortedset{sortedset: treemap.NewWithIntComparator()}
}
