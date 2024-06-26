package datastruct

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List test", func() {

	list := NewList()

	It("Size", func() {
		size := list.Size()
		Expect(size).To(Equal(int64(0)))

		list.PushLeft("1")
		Expect(list.Size()).To(Equal(int64(1)))
		list.PushRight("2")
		Expect(list.Size()).To(Equal(int64(2)))
	})

	It("PushLeft and PopLeft", func() {
		element1 := "1"
		element2 := "2"
		element3 := "3"
		list.PushLeft(element1)
		list.PushLeft(element2)
		list.PushLeft(element3)

		left, exist := list.PopLeft()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element3))

		left, exist = list.PopLeft()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element2))

		left, exist = list.PopLeft()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element1))
	})

	It("PushRight and PopRight", func() {
		element1 := "1"
		element2 := "2"
		element3 := "3"
		list.PushRight(element1)
		list.PushRight(element2)
		list.PushRight(element3)

		left, exist := list.PopRight()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element3))

		left, exist = list.PopRight()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element2))

		left, exist = list.PopRight()
		Expect(exist).To(BeTrue())
		Expect(left).To(Equal(element1))
	})
})
