package datastruct

import "math/rand"

type SkipListEntry struct {
	key   string
	value any

	left  *SkipListEntry
	right *SkipListEntry
	up    *SkipListEntry
	down  *SkipListEntry
}

// SkipList
//
// Reference: https://www.cs.emory.edu/~cheung/Courses/253/Syllabus/Map/skip-list-impl.html
type SkipList struct {
	head *SkipListEntry
	tail *SkipListEntry

	size  int
	level int
}

const maxLevel = 48

func NewSkipList() *SkipList {
	/*
			         +----+           +----+
			         |    +---------->|    |
		  level 0    |head|           |tail|
			         |    |<----------+    |
			         +----+           +----+
	*/
	head := &SkipListEntry{
		key:   "",
		value: nil,
	}
	tail := &SkipListEntry{
		key:   "\uffff",
		value: nil,
	}
	head.right = tail
	tail.left = head
	return &SkipList{
		head: head,
		tail: tail,
	}
}

func (skiplist *SkipList) findEntry(key string) (p *SkipListEntry) {
	p = skiplist.head
	for {
		for p.right != skiplist.tail && p.right.key <= key {
			p = p.right
		}

		// Go down next level if exists
		if p.down != nil {
			p = p.down
		} else {
			break
		}
	}
	if p == skiplist.tail {
		panic("tail.key should be empty")
	}
	return p
}

func (skiplist *SkipList) Put(key string, value any) (oldValue any) {
	// p.key <= key
	p := skiplist.findEntry(key)
	// If key exists, update value, then return old value. Done!
	if p.key == key {
		oldValue = p.value
		p.value = value
		return oldValue
	}

	// Create new entry
	newEntry := &SkipListEntry{
		key:   key,
		value: value,
	}

	/*
		        +------+            +------+            +------+
		level 1 |      |------------|      |------------|      |
		        +------+            +------+            +------+
		           |                   |                   |
		        +------+  +------+  +------+  +------+  +------+
		level 0 |      |--|      |--|      |--|      |--|      |
		        +------+  +------+  +------+  +------+  +------+
	*/
	// At level 0
	// p is the previous entry and p.key < key
	newEntry.left = p
	newEntry.right = p.right
	p.right.left = newEntry
	p.right = newEntry

	// toss a coin to decide whether to move up
	currentLevel := 0
	for rand.Intn(2) == 0 && currentLevel < maxLevel {
		currentLevel++
		if currentLevel > skiplist.level {
			// Create new level
			newLevelHead := &SkipListEntry{
				key:   "",
				value: nil,
			}
			newLevelTail := &SkipListEntry{
				key:   "\uffff",
				value: nil,
			}
			// Connect each other
			newLevelHead.right = newLevelTail
			newLevelTail.left = newLevelHead

			// Connect to the lower level head and tail
			newLevelHead.down = skiplist.head
			newLevelTail.down = skiplist.tail
			skiplist.head.up = newLevelHead
			skiplist.tail.up = newLevelTail

			// Update head and tail
			skiplist.head = newLevelHead
			skiplist.tail = newLevelTail

			skiplist.level++
		}

		// Reserve current newEntry (p.right)
		newEntryAtCurrentLevel := p.right
		// Go left until find an entry which can go up
		for p != skiplist.head && p.up == nil {
			p = p.left
		}
		p = p.up

		// Create new entry at upper level
		newEntryAtUpperLevel := &SkipListEntry{
			key:   key,
			value: nil, // do not need
		}
		newEntryAtUpperLevel.left = p
		newEntryAtUpperLevel.right = p.right
		newEntryAtUpperLevel.down = newEntryAtCurrentLevel

		p.right.left = newEntryAtUpperLevel
		p.right = newEntryAtUpperLevel

		newEntryAtCurrentLevel.up = newEntryAtUpperLevel
	}
	skiplist.size++
	return nil // no old value
}

func (skiplist *SkipList) Remove(key string) (value any) {
	p := skiplist.findEntry(key)
	if p.key != key {
		return nil
	}
	value = p.value
	for p != nil {
		p.left.right = p.right
		p.right.left = p.left
		// 清理指针
		p.left = nil
		p.right = nil
		p.up = nil
		p.down = nil
		p = p.up
	}
	skiplist.size--

	// 检查是否需要减少层数
	for skiplist.level > 0 && skiplist.head.right == skiplist.tail {
		// 最高层只剩 head 和 tail，移除该层
		skiplist.head = skiplist.head.down
		skiplist.tail = skiplist.tail.down
		skiplist.head.up = nil
		skiplist.tail.up = nil
		skiplist.level--
	}
	return value
}
func (skiplist *SkipList) Get(key string) (v any) {
	p := skiplist.findEntry(key)
	if p.key == key {
		return p.value
	}
	return nil
}
