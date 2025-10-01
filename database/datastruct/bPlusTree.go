package datastruct

import (
	"fmt"
	"github.com/noexcs/redis-go/log"
	"slices"
	"sync"
)

// MaxChildren 内部节点的孩子的个数，
// 又称 Order (used by Knuth's definition)
// MaxChildren + 1 >= 2 * MinChildren : 确保分裂后的两个Node 满足 MinChildren
const MaxChildren = 7

// MinChildren 内部节点的孩子的最少个数，
// 又称 Degree (used in the definition in Cormen et al. in Introduction to Algorithms (CLRS))
// MinChildren + (MinChildren - 1) <= MaxChildren : 确保合并后的可能的最大孩子数量不超过 MaxChildren
const MinChildren = 4

// MaxKeys 叶节点中的Key的最多个数
// MaxKeys + 1 >= 2 * MinKeys : 确保分裂后的两个 Node 满足 MinKeys
const MaxKeys = 3

// MinKeys 叶节点中的Key的最少个数
// MinKey + (MinKeys - 1) <= MaxKeys : 确保合并后的可能的最大孩子数量不超过 MaxKeys
const MinKeys = 1

type BPTKey interface {
	int | string | VolatileKey | interface{}
}

type UnionNode[T BPTKey] struct {
	parent   *UnionNode[T]
	leftPtr  *UnionNode[T]
	rightPtr *UnionNode[T]
	isLeaf   bool

	// internalNode field
	keys     []T
	children []*UnionNode[T]

	// leafNode field
	kvPairs []*KVPair[T]
}

func newUnionNode[T BPTKey](parent *UnionNode[T], leftPtr *UnionNode[T], rightPtr *UnionNode[T], isLeaf bool) *UnionNode[T] {
	node := &UnionNode[T]{
		parent:   parent,
		leftPtr:  leftPtr,
		rightPtr: rightPtr,
		isLeaf:   isLeaf,
	}
	if isLeaf {
		// 由于实现是每次添加后检查是否大于 MaxKeys 再进行分裂的，所以 kvPairs 最多能达到 MaxKeys + 1
		node.kvPairs = make([]*KVPair[T], 0, MaxKeys+1)
	} else {
		// 由于实现是每次添加后检查是否大于 MaxChildren 再进行分裂的，所以 children 最多能达到 MaxChildren + 1
		node.keys = make([]T, 0, MaxChildren)
		node.children = make([]*UnionNode[T], 0, MaxChildren+1)
	}
	return node
}

func (node *UnionNode[T]) insertKeyValue(t *BPlusTree[T], key T, value any, replaceKey bool) {
	idx, found := slices.BinarySearchFunc(node.kvPairs, key, func(kvPair *KVPair[T], k T) int {
		return CompareBPTKeys[T](kvPair.key, k)
	})
	if found {
		node.kvPairs[idx].value = value
		if replaceKey {
			node.kvPairs[idx].key = key
		}
	} else {
		node.kvPairs = append(node.kvPairs[:idx], append([]*KVPair[T]{{key, value}}, node.kvPairs[idx:]...)...)
	}
	if len(node.kvPairs) > MaxKeys {
		node.split(t)
	}
}

func (node *UnionNode[T]) insertNode(t *BPlusTree[T], key T, childNode *UnionNode[T]) {
	insertKeyIdx, _ := slices.BinarySearchFunc(node.keys, key, func(keyInNode, k T) int {
		return CompareBPTKeys[T](keyInNode, k)
	})
	node.keys = append(node.keys[:insertKeyIdx], append([]T{key}, node.keys[insertKeyIdx:]...)...)
	node.children = append(node.children[:insertKeyIdx+1], append([]*UnionNode[T]{childNode}, node.children[insertKeyIdx+1:]...)...)
	if len(node.children) > MaxChildren {
		node.split(t)
	}
}

func (node *UnionNode[T]) removeChild(t *BPlusTree[T], childNode *UnionNode[T]) {
	for idx, childPtr := range node.children {
		if childPtr == childNode {
			if idx-1 == -1 {
				node.keys = node.keys[1:]
			} else {
				node.keys = append(node.keys[:idx-1], node.keys[idx:]...)
			}
			node.children = append(node.children[:idx], node.children[idx+1:]...)
			if node.parent == t.root && len(node.children) == 1 {
				t.root = node.children[0]
				t.root.parent = nil
			}
			break
		}
	}
	if node != t.root && len(node.children) < MinChildren {
		if !node.borrow() {
			node.merge(t)
		}
	}
}

func (node *UnionNode[T]) removeKey(t *BPlusTree[T], key T) bool {
	idx, found := slices.BinarySearchFunc(node.kvPairs, key, func(kvPair *KVPair[T], k T) int {
		return CompareBPTKeys[T](kvPair.key, k)
	})
	if found {
		node.kvPairs = append(node.kvPairs[:idx], node.kvPairs[idx+1:]...)
	}
	if node != t.root && len(node.kvPairs) < MinKeys {
		if !node.borrow() {
			node.merge(t)
		}
	}
	return found
}

func (node *UnionNode[T]) split(t *BPlusTree[T]) {
	// 新节点在父节点的右边
	newNode := newUnionNode(node.parent, node, node.rightPtr, node.isLeaf)
	var midKey T
	if node.isLeaf {
		midKeyIdx := len(node.kvPairs) / 2
		midKey = node.kvPairs[midKeyIdx].key
		newNode.kvPairs = append(newNode.kvPairs, node.kvPairs[midKeyIdx:]...)
		node.kvPairs = node.kvPairs[:midKeyIdx]
	} else {
		// midKeyIdx 处的key移到父节点中，对应的child指针作为新节点左边的指针
		midKeyIdx := len(node.keys) / 2
		midKey = node.keys[midKeyIdx]

		// 新的节点只包含 midKeyIdx 后面的key
		newNode.keys = append(newNode.keys, node.keys[midKeyIdx+1:]...)
		newNode.children = append(newNode.children, node.children[midKeyIdx+1:]...)
		// 更改新节点的孩子的父节点为新节点
		for i := 0; i < len(newNode.children); i++ {
			newNode.children[i].parent = newNode
		}

		node.keys = node.keys[:midKeyIdx]
		node.children = node.children[:midKeyIdx+1]

		// 原本的 ChildPtr 不再是同一个父节点，更改对应的指针为 nil
		if !node.children[len(node.children)-1].isLeaf {
			node.children[len(node.children)-1].rightPtr = nil
			newNode.children[0].leftPtr = nil
		}
	}
	if node.rightPtr != nil {
		node.rightPtr.leftPtr = newNode
	}
	node.rightPtr = newNode
	if node.parent == nil {
		parent := newUnionNode[T](nil, nil, nil, false)
		parent.keys = append(parent.keys, midKey)
		parent.children = append(parent.children, node, newNode)

		node.parent = parent
		newNode.parent = parent
		t.root = parent
	} else {
		node.parent.insertNode(t, midKey, newNode)
	}
}

func (node *UnionNode[T]) borrow() bool {
	parent := node.parent
	leftNodePtr := node.leftPtr
	rightNodePtr := node.rightPtr
	idxInParent := getIdxInParent(node)
	if node.isLeaf {
		if borrowable[T](leftNodePtr, node) {
			// 先从左边借
			node.kvPairs = append([]*KVPair[T]{leftNodePtr.kvPairs[len(leftNodePtr.kvPairs)-1]}, node.kvPairs...)
			leftNodePtr.kvPairs = leftNodePtr.kvPairs[:len(leftNodePtr.kvPairs)-1]
			// 更新父节点指针对应的 key 的大小
			parent.keys[idxInParent-1] = node.kvPairs[0].key
			return true
		} else if borrowable(rightNodePtr, node) {
			// 再从右边借
			node.kvPairs = append(node.kvPairs, rightNodePtr.kvPairs[0])
			rightNodePtr.kvPairs = rightNodePtr.kvPairs[1:]
			// 更新父节点指针对应的key的大小
			parent.keys[idxInParent] = rightNodePtr.kvPairs[0].key
			return true
		}
	} else {
		if borrowable(leftNodePtr, node) {
			// 向左借（右旋）
			midKey := parent.keys[idxInParent-1]
			// 更改父节点 key
			parent.keys[idxInParent-1] = leftNodePtr.keys[len(leftNodePtr.keys)-1]

			// 更改该节点 key children 属性
			node.keys = append([]T{midKey}, node.keys...)
			node.children = append([]*UnionNode[T]{leftNodePtr.children[len(leftNodePtr.children)-1]}, node.children...)

			// 更改左节点 key children 属性
			leftNodePtr.keys = leftNodePtr.keys[:len(leftNodePtr.keys)-1]
			leftNodePtr.children = leftNodePtr.children[:len(leftNodePtr.children)-1]

			// 更改孩子节点父节点指针，以及左右指针
			leftNodePtr.children[len(leftNodePtr.children)-1].rightPtr = nil
			node.children[0].parent = node
			node.children[0].leftPtr = nil
			node.children[0].rightPtr = node.children[1]
			node.children[1].leftPtr = node.children[0]

			return true
		} else if borrowable(rightNodePtr, node) {
			// 向右借（左旋）
			midKey := parent.keys[idxInParent]
			// 更改父节点 key
			parent.keys[idxInParent] = rightNodePtr.keys[0]

			// 更改该节点 key children 属性
			node.keys = append(node.keys, midKey)
			node.children = append(node.children, rightNodePtr.children[0])

			// 更改右节点 key children 属性
			rightNodePtr.keys = rightNodePtr.keys[1:]
			rightNodePtr.children = rightNodePtr.children[1:]

			// 更改孩子节点父节点指针，以及左右指针
			rightNodePtr.children[0].leftPtr = nil
			node.children[len(node.children)-1].parent = node
			node.children[len(node.children)-1].rightPtr = nil
			node.children[len(node.children)-1].leftPtr = node.children[len(node.children)-2]
			node.children[len(node.children)-2].rightPtr = node.children[len(node.children)-1]

			return true
		}
	}
	return false
}

func (node *UnionNode[T]) merge(t *BPlusTree[T]) {
	parent := node.parent
	leftNodePtr := node.leftPtr
	rightNodePtr := node.rightPtr

	if mergeable(leftNodePtr) {
		// 这里的合并到左边
		rightNodePtr = node
		if !node.isLeaf {
			leftNodePtr.keys = append(leftNodePtr.keys, parent.keys[getIdxInParent(rightNodePtr)-1])
		}
	} else if mergeable(rightNodePtr) {
		// 右边的合并到这里
		leftNodePtr = node
		if !node.isLeaf {
			leftNodePtr.keys = append(leftNodePtr.keys, parent.keys[getIdxInParent(leftNodePtr)])
		}
	}

	if node.isLeaf {
		leftNodePtr.kvPairs = append(leftNodePtr.kvPairs, rightNodePtr.kvPairs...)
	} else {
		for _, childPtr := range rightNodePtr.children {
			childPtr.parent = leftNodePtr
		}
		leftNodePtr.keys = append(leftNodePtr.keys, rightNodePtr.keys...)

		leftNodePtr.children[len(leftNodePtr.children)-1].rightPtr = rightNodePtr.children[0]
		rightNodePtr.children[0].leftPtr = leftNodePtr.children[len(leftNodePtr.children)-1]
		leftNodePtr.children = append(leftNodePtr.children, rightNodePtr.children...)
	}

	leftNodePtr.rightPtr = rightNodePtr.rightPtr
	if rightNodePtr.rightPtr != nil {
		rightNodePtr.rightPtr.leftPtr = leftNodePtr
	}
	parent.removeChild(t, rightNodePtr)
}

func mergeable[T BPTKey](node *UnionNode[T]) bool {
	if node == nil {
		return false
	}
	if node.isLeaf {
		return len(node.kvPairs) == MinKeys
	} else {
		return len(node.children) == MinChildren
	}
}

func borrowable[T BPTKey](node *UnionNode[T], borrower *UnionNode[T]) bool {
	if node == nil || borrower == nil {
		return false
	}
	if node.isLeaf {
		return node.parent == borrower.parent && len(node.kvPairs) > MinKeys
	} else {
		return len(node.children) > MinChildren
	}
}

func getIdxInParent[T BPTKey](node *UnionNode[T]) int {
	if node.parent == nil {
		return -1
	}
	for idx, childPtr := range node.parent.children {
		if childPtr == node {
			return idx
		}
	}
	return -1
}

type KVPair[T BPTKey] struct {
	key   T
	value any
}

type BPlusTree[T BPTKey] struct {
	root          *UnionNode[T]
	firstLeafNode *UnionNode[T]
	mutex         sync.RWMutex
}

func MakeBPlusTree[T BPTKey]() *BPlusTree[T] {
	return &BPlusTree[T]{}
}

func (t *BPlusTree[T]) Insert(key T, value any, replaceKey bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.findLeafNode(key).insertKeyValue(t, key, value, replaceKey)
}

func (t *BPlusTree[T]) Delete(key T) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.findLeafNode(key).removeKey(t, key)
}

func (t *BPlusTree[T]) Clear() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.clearTree(t.root)
	t.root = nil
	t.firstLeafNode = nil
}

func (t *BPlusTree[T]) clearTree(node *UnionNode[T]) {
	if node == nil {
		return
	}

	// 清空子节点
	if node.isLeaf {
		for i := range node.kvPairs {
			node.kvPairs[i] = nil
		}
	} else {
		for i, child := range node.children {
			t.clearTree(child)
			node.children[i] = nil
		}
	}

	// 将节点数组置为nil，避免内存泄漏
	node.parent = nil
	node.leftPtr = nil
	node.rightPtr = nil
	node.keys = nil
	node.children = nil
	node.kvPairs = nil
}

func (t *BPlusTree[T]) Get(key T) (k T, v any, e bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	leafNode := t.findLeafNode(key)
	idx, found := slices.BinarySearchFunc(leafNode.kvPairs, key, func(kvPair *KVPair[T], k T) int {
		return CompareBPTKeys[T](kvPair.key, k)
	})
	if found {
		return leafNode.kvPairs[idx].key, leafNode.kvPairs[idx].value, true
	}
	return k, nil, false
}

// 返回key所在的leafNode, 或者应该插入的leafNode
func (t *BPlusTree[T]) findLeafNode(key T) (target *UnionNode[T]) {
	if t.root == nil {
		t.firstLeafNode = newUnionNode[T](nil, nil, nil, true)
		t.root = t.firstLeafNode
		return t.root
	}
	if t.root.isLeaf {
		return t.root
	}
	cursor := t.root
	for {
		idx, found := slices.BinarySearchFunc(cursor.keys, key, func(k1 T, k2 T) int {
			return CompareBPTKeys[T](k1, k2)
		})
		if found {
			cursor = cursor.children[idx+1]
		} else {
			cursor = cursor.children[idx]
		}
		if cursor.isLeaf {
			return cursor
		}
	}
}

func (t *BPlusTree[T]) Iterator() *Iterator[T] {
	t.mutex.RLock()
	return &Iterator[T]{
		node: t.firstLeafNode,
		idx:  0,
		tree: t,
	}
}

type Iterator[T BPTKey] struct {
	node      *UnionNode[T]
	idx       int
	tree      *BPlusTree[T]
	rUnLocked bool
}

func (iter *Iterator[T]) Next() bool {
	return iter.node != nil && iter.idx < len(iter.node.kvPairs)
}

func (iter *Iterator[T]) Value() (k T, v any) {
	if iter.rUnLocked {
		if iter.Next() {
			panic("The Iterator is unavailable. Discarded ahead of time.")
		} else {
			panic("The Iterator is unavailable. No more elements.")
		}
	}
	if !iter.Next() {
		panic("No more elements.")
	}
	kvPair := iter.node.kvPairs[iter.idx]
	if iter.idx+1 < len(iter.node.kvPairs) {
		iter.idx += 1
	} else {
		iter.node = iter.node.rightPtr
		iter.idx = 0
	}
	// 如果没有下一个（迭代结束了），且还没解锁读锁，解锁读锁
	if !iter.Next() && !iter.rUnLocked {
		iter.tree.mutex.RUnlock()
		iter.rUnLocked = true
	}
	return kvPair.key, kvPair.value
}

// Discard 不再使用该迭代器，提前解锁读锁
func (iter *Iterator[T]) Discard() {
	// 判断一下保证幂等性
	if !iter.rUnLocked {
		iter.node = nil
		iter.tree.mutex.RUnlock()
		iter.rUnLocked = true
		iter.tree = nil
	}
}

// CompareStrings 比较两个字符串，返回比较结果：
// 如果 str1 < str2，返回 -1；
// 如果 str1 == str2，返回 0；
// 如果 str1 > str2，返回 1。
func CompareStrings(str1, str2 string) int {
	if str1 == str2 {
		return 0
	} else if str1 < str2 {
		return -1
	} else {
		return 1
	}
}

func CompareBPTKeys[T BPTKey](key1 T, key2 T) int {
	var k1 string
	var k2 string
	switch v := any(key1).(type) {
	case string:
		k1 = v
	case *VolatileKey:
		k1 = v.Name
	default:
		log.WithLocation(fmt.Sprintf("Unsupported key type of %T", key1))
		return 0
	}
	switch v := any(key2).(type) {
	case string:
		k2 = v
	case *VolatileKey:
		k2 = v.Name
	default:
		log.WithLocation(fmt.Sprintf("Unsupported key type of %T", key1))
		return 0
	}
	return CompareStrings(k1, k2)
}
