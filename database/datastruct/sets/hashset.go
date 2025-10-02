package sets

import (
	"github.com/emirpasic/gods/sets/hashset"
)

// HashSet 使用hashset实现的Set
type HashSet struct {
	dict *hashset.Set
}

func NewHashSet() *HashSet {
	return &HashSet{dict: hashset.New()}
}

func (s *HashSet) Add(value string) {
	s.dict.Add(value)
}

func (s *HashSet) Members() []string {

	values := s.dict.Values()

	memberCnt := len(values)
	result := make([]string, memberCnt)

	for _, value := range values {
		result[memberCnt-1] = value.(string)
		memberCnt--
	}

	return result
}

func (s *HashSet) Remove(v string) {
	s.dict.Remove(v)
}

func (s *HashSet) Contains(v string) bool {
	return s.dict.Contains(v)
}

func (s *HashSet) Size() int {
	return s.dict.Size()
}
