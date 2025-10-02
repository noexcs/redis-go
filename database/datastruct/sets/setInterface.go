package sets

// Set 定义了Set数据结构的通用接口
type Set interface {
	// Add 向集合中添加元素
	Add(value string)

	// Members 返回集合中的所有元素
	Members() []string

	// Remove 从集合中移除元素
	Remove(v string)

	// Contains 检查元素是否在集合中
	Contains(v string) bool

	// Size 返回集合中元素的数量
	Size() int
}
