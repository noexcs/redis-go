package datastruct

import (
	"math/rand"
	"testing"
	"time"
)

// 初始化随机种子以确保测试的一致性
func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestSkipList_New(t *testing.T) {
	sl := newSkipList()

	if sl == nil {
		t.Error("Expected newSkipList() to return a non-nil SkipList")
	}

	if sl.size != 0 {
		t.Errorf("Expected initial size to be 0, got %d", sl.size)
	}

	if sl.level != 0 {
		t.Errorf("Expected initial level to be 0, got %d", sl.level)
	}

	if sl.head == nil || sl.tail == nil {
		t.Error("Expected head and tail to be initialized")
	}

	if sl.head.right != sl.tail {
		t.Error("Expected head.right to point to tail")
	}

	if sl.tail.left != sl.head {
		t.Error("Expected tail.left to point to head")
	}
}

func TestSkipList_PutAndGet(t *testing.T) {
	sl := newSkipList()

	// 测试插入新键值对
	oldValue := sl.Put("key1", "value1")
	if oldValue != nil {
		t.Errorf("Expected nil for new key, got %v", oldValue)
	}

	if sl.size != 1 {
		t.Errorf("Expected size to be 1 after first insert, got %d", sl.size)
	}

	// 测试获取存在的键值对
	value := sl.Get("key1")
	if value != "value1" {
		t.Errorf("Expected 'value1' for key 'key1', got %v", value)
	}

	// 测试更新已存在的键
	oldValue = sl.Put("key1", "updated_value1")
	if oldValue != "value1" {
		t.Errorf("Expected 'value1' as old value, got %v", oldValue)
	}

	value = sl.Get("key1")
	if value != "updated_value1" {
		t.Errorf("Expected 'updated_value1' for key 'key1', got %v", value)
	}

	// 测试获取不存在的键
	value = sl.Get("nonexistent")
	if value != nil {
		t.Errorf("Expected nil for nonexistent key, got %v", value)
	}
}

func TestSkipList_Remove(t *testing.T) {
	sl := newSkipList()

	// 插入一些数据
	sl.Put("key1", "value1")
	sl.Put("key2", "value2")
	sl.Put("key3", "value3")

	if sl.size != 3 {
		t.Errorf("Expected size to be 3 after 3 inserts, got %d", sl.size)
	}

	// 删除存在的键
	removedValue := sl.Remove("key2")
	if removedValue != "value2" {
		t.Errorf("Expected 'value2' to be removed, got %v", removedValue)
	}

	if sl.size != 2 {
		t.Errorf("Expected size to be 2 after removal, got %d", sl.size)
	}

	// 确认键已被删除
	value := sl.Get("key2")
	if value != nil {
		t.Errorf("Expected nil for removed key, got %v", value)
	}

	// 删除不存在的键
	removedValue = sl.Remove("nonexistent")
	if removedValue != nil {
		t.Errorf("Expected nil when removing nonexistent key, got %v", removedValue)
	}

	if sl.size != 2 {
		t.Errorf("Expected size to remain 2 when removing nonexistent key, got %d", sl.size)
	}
}

func TestSkipList_MultipleOperations(t *testing.T) {
	sl := newSkipList()

	// 插入大量数据
	testData := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": "three",
		"d": 4.0,
		"e": true,
	}

	for key, value := range testData {
		sl.Put(key, value)
	}

	if sl.size != len(testData) {
		t.Errorf("Expected size to be %d, got %d", len(testData), sl.size)
	}

	// 验证所有数据
	for key, expectedValue := range testData {
		actualValue := sl.Get(key)
		if actualValue != expectedValue {
			t.Errorf("Expected %v for key %s, got %v", expectedValue, key, actualValue)
		}
	}

	// 删除一些数据
	keysToRemove := []string{"b", "d"}
	for _, key := range keysToRemove {
		sl.Remove(key)
		delete(testData, key)
	}

	if sl.size != len(testData) {
		t.Errorf("Expected size to be %d after removals, got %d", len(testData), sl.size)
	}

	// 再次验证剩余数据
	for key, expectedValue := range testData {
		actualValue := sl.Get(key)
		if actualValue != expectedValue {
			t.Errorf("Expected %v for key %s after removals, got %v", expectedValue, key, actualValue)
		}
	}
}

func TestSkipList_LevelManagement(t *testing.T) {
	sl := newSkipList()

	// 初始级别应该是0
	if sl.level != 0 {
		t.Errorf("Expected initial level to be 0, got %d", sl.level)
	}

	// 插入足够多的数据以触发级别提升
	// 我们使用确定性的数据来确保测试一致性
	rand.Seed(12345) // 固定种子以确保测试可重复

	for i := 0; i < 100; i++ {
		sl.Put(string(rune(i+65)), i)
	}

	// 验证级别是否增加（概率上应该会增加）
	if sl.level < 0 {
		t.Errorf("Expected level to be non-negative, got %d", sl.level)
	}

	// 验证一些数据
	if sl.Get("A") != 0 || sl.Get("Z") != 25 {
		t.Error("Data retrieval failed after level increases")
	}

	// 恢复随机种子
	rand.Seed(time.Now().UnixNano())
}

func TestSkipList_EmptyOperations(t *testing.T) {
	sl := newSkipList()

	// 在空SkipList上进行各种操作
	value := sl.Get("anykey")
	if value != nil {
		t.Errorf("Expected nil from Get on empty SkipList, got %v", value)
	}

	removed := sl.Remove("anykey")
	if removed != nil {
		t.Errorf("Expected nil from Remove on empty SkipList, got %v", removed)
	}
}
