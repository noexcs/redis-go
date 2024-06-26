package datastruct

import (
	"github.com/emirpasic/gods/maps/treemap"
)

// Sortedset https://redis.io/docs/data-types/sorted-sets/
type Sortedset struct {
	sortedset *treemap.Map
}

func (s *Sortedset) Add(score int64, elem string) {
	s.sortedset.Put(score, elem)
}

func NewSortedset() *Sortedset {
	return &Sortedset{sortedset: treemap.NewWithIntComparator()}
}
