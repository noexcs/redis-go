package benchmark

import (
	"github.com/noexcs/redis-go/database/datastruct/sets"
	"math/rand"
	"strconv"
	"testing"
)

// goos: windows
// goarch: amd64
// pkg: github.com/noexcs/redis-go/benchmark
// cpu: 13th Gen Intel(R) Core(TM) i5-13600KF

// BenchmarkMapSetAdd
// BenchmarkMapSetAdd-20          	 5631198	       239.6 ns/op
// BenchmarkMapSetContains
// BenchmarkMapSetContains-20     	27238059	        44.07 ns/op
// BenchmarkMapSetRemove
// BenchmarkMapSetRemove-20       	39897860	        28.96 ns/op

// BenchmarkHashSetAdd
// BenchmarkHashSetAdd-20         	 4110613	       396.7 ns/op
// BenchmarkHashSetContains
// BenchmarkHashSetContains-20    	22635444	        51.17 ns/op
// BenchmarkHashSetRemove
// BenchmarkHashSetRemove-20      	34055402	        33.46 ns/op
func BenchmarkMapSetAdd(b *testing.B) {
	s := sets.NewMapSet()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}
}

func BenchmarkMapSetContains(b *testing.B) {
	s := sets.NewMapSet()
	// Pre-populate the set
	for i := 0; i < 10000; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(rand.Intn(10000))
		_ = s.Contains(value)
	}
}

func BenchmarkMapSetRemove(b *testing.B) {
	s := sets.NewMapSet()
	// Pre-populate the set
	for i := 0; i < 10000; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(rand.Intn(10000))
		s.Remove(value)
	}
}

func BenchmarkHashSetAdd(b *testing.B) {
	s := sets.NewHashSet()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}
}

func BenchmarkHashSetContains(b *testing.B) {
	s := sets.NewHashSet()
	// Pre-populate the set
	for i := 0; i < 10000; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(rand.Intn(10000))
		_ = s.Contains(value)
	}
}

func BenchmarkHashSetRemove(b *testing.B) {
	s := sets.NewHashSet()
	// Pre-populate the set
	for i := 0; i < 10000; i++ {
		value := "value" + strconv.Itoa(i)
		s.Add(value)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		value := "value" + strconv.Itoa(rand.Intn(10000))
		s.Remove(value)
	}
}
