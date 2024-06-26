package command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("hash", func() {

	key := "hash_example_key"

	f1, v1 := "f1", "v1"
	f2, v2 := "f2", "v2"
	f3, v3 := "f3", "v3"
	f4, _ := "f4", "v4"

	It("hset", func() {
		result, err := client.HSet(ctx, key, f1, v1, f2, v2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))

		result, err = client.HSet(ctx, key, f3, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))
	})

	It("hsetnx", func() {
		key := "hsetnx_key"
		field1 := "field1"
		value1 := "value1"
		field2 := "field2"
		value2 := "value2"
		client.HMSet(ctx, key, field1, value1)

		result, err := client.HSetNX(ctx, key, field1, value2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(false))

		result, err = client.HSetNX(ctx, key, field2, value2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(true))

		s, err := client.HGet(ctx, key, field1).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(value1))

		s, err = client.HGet(ctx, key, field2).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal(value2))
	})

	It("hget", func() {
		result, err := client.HGet(ctx, key, f1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(v1))

		result, err = client.HGet(ctx, key, f4).Result()
		Expect(err).NotTo(BeNil())
		Expect(result).To(Equal(""))
	})

	It("hmset", func() {
		key := "hmset_key"
		field1 := "field1"
		value1 := "value1"
		field2 := "field2"
		value2 := "value2"
		client.HMSet(ctx, key, field1, value1, field2, value2)

		result, err := client.HGet(ctx, key, field1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value1))

		result, err = client.HGet(ctx, key, field2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value2))
	})

	It("hmget", func() {
		key := "hmget_key"
		field1 := "field1"
		value1 := "value1"
		field2 := "field2"
		value2 := "value2"
		client.HMSet(ctx, key, field1, value1, field2, value2)

		result, err := client.HMGet(ctx, key, field1, field2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]any{value1, value2}))
	})

	It("hdel", func() {
		result, err := client.HDel(ctx, key, f1, f2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))

		value, err := client.HGet(ctx, key, f1).Result()
		Expect(err).NotTo(BeNil())
		Expect(value).To(Equal(""))

		value, err = client.HGet(ctx, key, f2).Result()
		Expect(err).NotTo(BeNil())
		Expect(value).To(Equal(""))
	})

	It("hexists", func() {
		result, err := client.HExists(ctx, key, f1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(false))

		result, err = client.HExists(ctx, key, f3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(true))
	})

	It("hgetall", func() {
		result, err := client.HGetAll(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(map[string]string{f3: v3}))
	})

	It("hincrby", func() {
		client.HSet(ctx, key, f1, "1")

		result, err := client.HIncrBy(ctx, key, f1, 1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("HINCBYFLOAT", func() {
		key := "HINCBYFLOAT_key1"
		f1 := "HINCBYFLOAT_f1"

		client.HSet(ctx, key, f1, "10.5")

		result, err := client.HIncrByFloat(ctx, key, f1, 1.1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(11.6))
	})

	It("hkeys", func() {
		result, err := client.HKeys(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{f3, f1}))
	})

	It("hlen", func() {
		result, err := client.HLen(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	// unimplemented
	PIt("hstrlen", func() {
	})

	It("hvals", func() {
		result, err := client.HVals(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{v3, "2"}))
	})

	It("HRANDFIELD", func() {
		key := "HRANDFIELD_key"

		set := make(map[string]string)
		for i := 0; i < 10; i++ {
			value := strconv.Itoa(i)
			field := "field" + value
			set[field] = value
			client.HSet(ctx, key, value)
		}

		result, err := client.HRandField(ctx, key, 2, true).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Satisfy(func(r []string) bool {
			if len(r)%2 != 0 {
				return false
			}
			for i := 0; i < len(r); i += 2 {
				if _, ok := set[r[i]]; !ok {
					return false
				}
				if i+1 < len(r) {
					return false
				}
				if r[i+1] != set[r[i]] {
					return false
				}
			}
			return true
		}))
	})
})
