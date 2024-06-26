package command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("string", func() {

	It("SET and GET", func() {
		key := "SET_GET_key"
		value := "Hello, Redis"

		result, err := client.Set(ctx, key, value, 0).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("OK"))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(value))
	})

	It("SET NX", func() {
		key := "SETNX_key"
		value := "Hello, Redis"

		result, err := client.Set(ctx, key, value, 0).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("OK"))

		b, err := client.SetNX(ctx, key, "Hello, World", 0).Result()
		Expect(err).To(BeNil())
		Expect(b).To(Equal(false))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(value))
	})

	It("SET XX", func() {
		key := "SETXX_key"
		value := "Hello, Redis"

		result, err := client.Set(ctx, key, value, 0).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("OK"))

		newValue := "Hello, World"
		b, err := client.SetXX(ctx, key, newValue, 0).Result()
		Expect(err).To(BeNil())
		Expect(b).To(Equal(true))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(newValue))
	})

	It("SET EX", func() {
		key := "SETEX_key"

		value := "Hello, Redis"
		client.SetEX(ctx, key, value, 5*time.Second)

		time.Sleep(2 * time.Second)
		result, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value))

		time.Sleep(3 * time.Second)
		s, err := client.Get(ctx, key).Result()
		Expect(err).NotTo(BeNil())
		Expect(s).To(Equal(""))
	})

	It("SETRANGE", func() {
		key := "SETRANGE_key1"
		value := "Hello Redis"

		client.Set(ctx, key, value, 0)

		// 设置字符串的一部分
		client.SetRange(ctx, key, 6, "World").Result()

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal("Hello World"))
	})

	Context("GETRANGE", func() {
		It("GETRANGE: in range", func() {
			key := "GETRANGE_key1"
			value := "Hello, Redis"

			client.Set(ctx, key, value, 0)

			val, err := client.GetRange(ctx, key, 0, 4).Result()
			Expect(err).To(BeNil())
			Expect(val).To(Equal(value[0 : 4+1]))
		})

		It("GETRANGE: out of range", func() {
			key := "GETRANGE_key2"
			value := "Hello, Redis"

			client.Set(ctx, key, value, 0)

			val, err := client.GetRange(ctx, key, 0, 100).Result()
			Expect(err).To(BeNil())
			Expect(val).To(Equal(value))
		})
	})

	It("GETSET", func() {
		key := "GETSET_key"
		value1 := "Hello"
		value2 := "Redis"

		client.Set(ctx, key, value1, 0)
		result, err := client.GetSet(ctx, key, value2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value1))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(value2))
	})

	It("MGET", func() {
		key1 := "MGET_key1"
		key2 := "MGET_key2"
		value1 := "Hello"
		value2 := "Redis"

		client.Set(ctx, key1, value1, 0)
		client.Set(ctx, key2, value2, 0)

		result, err := client.MGet(ctx, key1, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]any{value1, value2}))
	})

	It("MSET", func() {
		key1 := "MSET_key1"
		key2 := "MSET_key2"
		value1 := "Hello"
		value2 := "Redis"

		client.MSet(ctx, key1, value1, key2, value2)

		result, err := client.MGet(ctx, key1, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]any{value1, value2}))
	})

	It("MSETNX", func() {
		key1 := "MSETNX_key1"
		key2 := "MSETNX_key2"
		key3 := "MSETNX_key3"
		value1 := "Hello"
		value2 := "there"
		value3 := "world"

		result, err := client.MSetNX(ctx, key1, value1, key2, value2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		result, err = client.MSetNX(ctx, key2, "new", key3, value3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(BeFalse())

		i, err := client.MGet(ctx, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(i).To(Equal([]any{value1, value2, nil}))
	})

	It("STRLEN", func() {
		key := "STRLEN_key"
		client.Set(ctx, key, "hello world", 0)

		result, err := client.StrLen(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(11)))

		result, err = client.StrLen(ctx, "STRLEN_nonexisting").Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(0)))
	})

	It("INCR", func() {
		key := "INCR_key"
		client.Set(ctx, key, "1", 0)

		result, err := client.Incr(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("INCRBY", func() {
		key := "INCRBY_key"
		client.Set(ctx, key, "10", 0)

		result, err := client.Incr(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(11)))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal("11"))
	})

	It("DECR", func() {
		key := "DECR_key"
		client.Set(ctx, key, "10", 0)

		result, err := client.Decr(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(9)))
	})

	It("DECRBY", func() {
		key := "DECRBY_key"
		client.Set(ctx, key, "10", 0)

		result, err := client.DecrBy(ctx, key, 3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(7)))
	})

	It("INCRBYFLOAT", func() {
		key := "INCRBYFLOAT_key"
		client.Set(ctx, key, "10.1", 0)

		result, err := client.IncrByFloat(ctx, key, 3.2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(13.3))
	})

	It("APPEND", func() {
		key := "APPEND_key"
		client.Set(ctx, key, "hello", 0)
		result, err := client.Append(ctx, key, " world").Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(11)))

		s, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(s))
	})

	It("GETBIT", func() {
		key := "GETBIT_key"
		client.SetBit(ctx, key, 7, 1)

		result, err := client.GetBit(ctx, key, 7).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))
	})

	Context("SETBIT", func() {
		It("SETBIT", func() {
			key := "SETBIT_key"
			result, err := client.SetBit(ctx, key, 7, 1).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(0)))

			result, err = client.SetBit(ctx, key, 7, 0).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))

			s, err := client.Get(ctx, key).Result()
			Expect(err).To(BeNil())
			Expect(s).To(Equal("\x00"))
		})
		It("bitmapsarestrings", func() {
			key := "SETBIT_key2"
			client.SetBit(ctx, key, 2, 1)
			client.SetBit(ctx, key, 3, 1)
			client.SetBit(ctx, key, 5, 1)
			client.SetBit(ctx, key, 10, 1)
			client.SetBit(ctx, key, 11, 1)
			client.SetBit(ctx, key, 14, 1)

			result, err := client.Get(ctx, key).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal("42"))
		})
	})
})
