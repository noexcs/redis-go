package datastruct

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("hashmap", func() {

	hashmap := NewHashmap()

	k1, v1 := "k1", "v1"
	k2 := "k2"

	It("Put", func() {
		hashmap.Put(k1, v1)
		contains := hashmap.Contains(k1)
		Expect(contains).To(Equal(true))
	})

	It("Get", func() {
		value, b := hashmap.Get(k1)
		Expect(b).To(Equal(true))
		Expect(value).To(Equal(v1))
	})

	It("Contains", func() {
		contains := hashmap.Contains(k2)
		Expect(contains).To(Equal(false))

		hashmap.Put(k1, v1)
		contains = hashmap.Contains(k1)
		Expect(contains).To(Equal(true))
	})
})
