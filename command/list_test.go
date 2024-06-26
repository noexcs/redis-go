package command_test

import (
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("list", func() {

	It("LPUSH and LPOP", func() {
		key := "LPUSH_LPOP_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		result, err := client.LPush(ctx, key, v1, v2, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(3)))

		value, err := client.LPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v3))

		value, err = client.LPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v2))

		value, err = client.LPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v1))
	})

	It("RPUSH and RPOP", func() {
		key := "RPUSH_RPOP_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		result, err := client.RPush(ctx, key, v1, v2, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(3)))

		value, err := client.RPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v3))

		value, err = client.RPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v2))

		value, err = client.RPop(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v1))
	})

	PIt("BLPOP and BRPOP", func() {

	})

	It("LINDEX", func() {
		key := "LINDEX_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.LPush(ctx, key, v1, v2, v3)

		value, err := client.LIndex(ctx, key, 0).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v3))

		value, err = client.LIndex(ctx, key, 1).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v2))

		value, err = client.LIndex(ctx, key, 2).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v1))
	})

	It("LLEN", func() {
		key := "LLEN_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.LPush(ctx, key, v1, v2, v3)
		result, err := client.LLen(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(3)))
	})

	It("LINSERT", func() {
		key := "LINSERT_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"
		v4 := "v4"

		client.RPush(ctx, key, v1, v2, v4)
		result, err := client.LInsert(ctx, key, "BEFORE", v4, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(4)))

		value, err := client.LIndex(ctx, key, 2).Result()
		Expect(err).To(BeNil())
		Expect(value).To(Equal(v3))
	})

	It("LRANGE", func() {
		key := "LRANGE_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"
		v4 := "v4"

		client.RPush(ctx, key, v1, v2, v3, v4)
		result, err := client.LRange(ctx, key, 1, 3).Result()
		Expect(err).To(BeNil())
		sort.Strings(result)
		Expect(result).To(Equal([]string{v2, v3, v4}))
	})

	It("LREM", func() {
		key := "LREM_key1"
		v1 := "v1"
		v2 := "v2"

		client.RPush(ctx, key, v1, v1, v2, v1)
		result, err := client.LRem(ctx, key, -2, v1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("LSET", func() {
		key := "LSET_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"
		v4 := "v4"
		v5 := "v5"

		client.RPush(ctx, key, v1, v2, v3)
		client.LSet(ctx, key, 0, v4)
		client.LSet(ctx, key, -2, v5)

		result, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{v4, v5, v3}))
	})

	It("LTRIM", func() {
		key := "LTRIM_key1"
		client.RPush(ctx, key, "v1", "v2", "v3")

		client.LTrim(ctx, key, 1, -1)
		result, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{"v2", "v3"}))
	})

	It("LPUSHX", func() {
		key := "LPUSHX_key"
		key2 := "LPUSHX_key2"

		client.LPush(ctx, key, "World")
		client.LPushX(ctx, key, "Hello")
		client.LPushX(ctx, key2, "Hello")

		result, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{"Hello", "World"}))

		strings, err := client.LRange(ctx, key2, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{}))
	})

	It("RPUSHX", func() {
		key := "RPUSHX_key"
		key2 := "RPUSHX_key2"

		client.RPush(ctx, key, "Hello")
		client.RPushX(ctx, key, "World")
		client.RPushX(ctx, key2, "World")

		result, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{"Hello", "World"}))

		strings, err := client.LRange(ctx, key2, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{}))
	})

	It("LMOVE", func() {
		key := "LMOVE_key"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.RPush(ctx, key, v1, v2, v3)
		key_dest := "LMOVE_key_dest"
		result, err := client.LMove(ctx, key, key_dest, "RIGHT", "LEFT").Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(v3))

		result, err = client.LMove(ctx, key, key_dest, "LEFT", "RIGHT").Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(v1))

		strings, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{v2}))

		strings, err = client.LRange(ctx, key_dest, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{v3, v1}))
	})

	It("LPOS", func() {
		key := "LPOS_key"

		client.RPush(ctx, key, "a", "b", "c", 1, 2, 3, "c", "c")
		result, err := client.LPos(ctx, key, "c", redis.LPosArgs{
			Rank:   2,
			MaxLen: 0,
		}).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(6)))

		result, err = client.LPos(ctx, key, "c", redis.LPosArgs{
			Rank:   -1,
			MaxLen: 0,
		}).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(7)))

		result, err = client.LPos(ctx, key, "c", redis.LPosArgs{
			Rank:   -1,
			MaxLen: 0,
		}).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(7)))
	})
})
