package benchmark

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/noexcs/redis-go/database/datastruct"
)

// goos: windows
// goarch: amd64
// pkg: github.com/noexcs/redis-go/benchmark
// cpu: 13th Gen Intel(R) Core(TM) i5-13600KF

// BenchmarkGoMapPut
// BenchmarkGoMapPut-20        	 4180278	       317.0 ns/op
// BenchmarkGoMapGet
// BenchmarkGoMapGet-20        	26782246	        45.79 ns/op

// BenchmarkHashmapPut
// BenchmarkHashmapPut-20      	 3274002	       349.2 ns/op
// BenchmarkHashmapGet
// BenchmarkHashmapGet-20      	22066165	        52.20 ns/op

// BenchmarkSkipListPut
// BenchmarkSkipListPut-20     	 4208809	       322.1 ns/op
// BenchmarkSkipListGet
// BenchmarkSkipListGet-20     	 5464201	       222.0 ns/op

// BenchmarkBPlusTreePut
// BenchmarkBPlusTreePut-20    	 2791926	       459.5 ns/op
// BenchmarkBPlusTreeGet
// BenchmarkBPlusTreeGet-20    	 4428097	       267.8 ns/op

func BenchmarkGoMapPut(b *testing.B) {
	m := make(map[string]string)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		m[key] = "value" + strconv.Itoa(i)
	}
}

func BenchmarkGoMapGet(b *testing.B) {
	m := make(map[string]string)
	// Pre-populate the map
	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(i)
		m[key] = "value" + strconv.Itoa(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(rand.Intn(10000))
		_ = m[key]
	}
}

func BenchmarkHashmapPut(b *testing.B) {
	h := datastruct.NewHashmap()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		h.Put(key, "value"+strconv.Itoa(i))
	}
}

func BenchmarkHashmapGet(b *testing.B) {
	h := datastruct.NewHashmap()
	// Pre-populate the hashmap
	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(i)
		h.Put(key, "value"+strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(rand.Intn(10000))
		_, _ = h.Get(key)
	}
}

func BenchmarkSkipListPut(b *testing.B) {
	sl := datastruct.NewSkipList()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		sl.Put(key, "value"+strconv.Itoa(i))
	}
}

func BenchmarkSkipListGet(b *testing.B) {
	sl := datastruct.NewSkipList()
	// Pre-populate the skiplist
	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(i)
		sl.Put(key, "value"+strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(rand.Intn(10000))
		_ = sl.Get(key)
	}
}

func BenchmarkBPlusTreePut(b *testing.B) {
	bt := datastruct.MakeBPlusTree[string]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		bt.Insert(key, "value"+strconv.Itoa(i), false)
	}
}

func BenchmarkBPlusTreeGet(b *testing.B) {
	bt := datastruct.MakeBPlusTree[string]()
	// Pre-populate the B+ tree
	for i := 0; i < 10000; i++ {
		key := "key" + strconv.Itoa(i)
		bt.Insert(key, "value"+strconv.Itoa(i), false)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(rand.Intn(10000))
		_, _, _ = bt.Get(key)
	}
}
