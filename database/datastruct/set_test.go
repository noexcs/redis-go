package datastruct

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set test", func() {
	It("Size", func() {
		set := NewSet()
		size := set.Size()
		Expect(size).To(Equal(0))
	})
})
