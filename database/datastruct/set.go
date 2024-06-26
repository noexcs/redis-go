package datastruct

import (
	"github.com/emirpasic/gods/sets/hashset"
)

type Set struct {
	dict *hashset.Set
}

func NewSet() *Set {
	return &Set{dict: hashset.New()}
}

func (s *Set) Add(value string) {
	s.dict.Add(value)
}

func (s *Set) Members() []string {

	values := s.dict.Values()

	memberCnt := len(values)
	result := make([]string, memberCnt)

	for _, value := range values {
		result[memberCnt-1] = value.(string)
		memberCnt--
	}

	return result
}

func (s *Set) Remove(v string) {
	s.dict.Remove(v)
}

func (s *Set) Contains(v string) bool {
	return s.dict.Contains(v)
}

func (s *Set) Size() int {
	return s.dict.Size()
}
