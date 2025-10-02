package sets

// MapSet 使用Go内置map实现的Set
type MapSet struct {
	dict map[string]bool
}

// NewMapSet 创建一个新的MapSet实例
func NewMapSet() *MapSet {
	return &MapSet{dict: make(map[string]bool)}
}

// Add 向集合中添加元素
func (s *MapSet) Add(value string) {
	s.dict[value] = true
}

// Members 返回集合中的所有元素
func (s *MapSet) Members() []string {
	result := make([]string, 0, len(s.dict))
	for k := range s.dict {
		result = append(result, k)
	}
	return result
}

// Remove 从集合中移除元素
func (s *MapSet) Remove(v string) {
	delete(s.dict, v)
}

// Contains 检查元素是否在集合中
func (s *MapSet) Contains(v string) bool {
	return s.dict[v]
}

// Size 返回集合中元素的数量
func (s *MapSet) Size() int {
	return len(s.dict)
}
