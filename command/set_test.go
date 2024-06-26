package command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("set", func() {

	It("sadd", func() {
		key := "sadd_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		result, err := client.SAdd(ctx, key, v1, v2, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(3)))
	})

	It("srem", func() {
		key := "srem_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.SAdd(ctx, key, v1, v2, v3)

		result, err := client.SRem(ctx, key, v2, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("smembers", func() {
		key := "smembers_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.SAdd(ctx, key, v1, v2, v3)

		result, err := client.SMembers(ctx, key).Result()
		Expect(err).To(BeNil())
		sort.Strings(result)
		expected := []string{v1, v2, v3}
		Expect(result).To(Equal(expected))
	})

	It("sismember", func() {
		key := "sismember_key1"
		v1 := "v1"
		v2 := "v2"
		v3 := "v3"

		client.SAdd(ctx, key, v1, v2)

		result, err := client.SIsMember(ctx, key, v1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(true))

		result, err = client.SIsMember(ctx, key, v3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(false))
	})

	It("scard", func() {
		key := "scard_key1"
		v1 := "v1"
		v2 := "v2"

		client.SAdd(ctx, key, v1, v2)

		result, err := client.SCard(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("SDIFF", func() {
		key1, a, b, c, d := "SDIFF_key1", "a", "b", "c", "d"
		key2 := "SDIFF_key2"
		key3 := "SDIFF_key3"

		client.SAdd(ctx, key1, a, b, c, d)
		client.SAdd(ctx, key2, c)
		client.SAdd(ctx, key3, b, d)

		result, err := client.SDiff(ctx, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{a}))

	})

	It("SDIFFSTORE", func() {
		key1, a, b, c, d := "SDIFFSTORE_key1", "a", "b", "c", "d"
		key2 := "SDIFFSTORE_key2"
		key3 := "SDIFFSTORE_key3"
		client.SAdd(ctx, key1, a, b, c, d)
		client.SAdd(ctx, key2, c)
		client.SAdd(ctx, key3, b, d)

		destination := "dest"

		result, err := client.SDiffStore(ctx, destination, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))

		s, err := client.SMembers(ctx, destination).Result()
		Expect(err).To(BeNil())
		Expect(s).To(Equal([]string{a}))
	})

	It("SINTER", func() {
		key1, a, b, c, d, e := "SINTER_key1", "a", "b", "c", "d", "e"
		key2 := "SINTER_key2"
		key3 := "SINTER_key3"
		client.SAdd(ctx, key1, a, b, c, d)
		client.SAdd(ctx, key2, c)
		client.SAdd(ctx, key3, a, c, e)

		result, err := client.SInter(ctx, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]string{c}))
	})

	It("SINTERSTORE", func() {
		key1, a, b, c, d, e := "SINTERSTORE_key1", "a", "b", "c", "d", "e"
		key2 := "SINTERSTORE_key2"
		client.SAdd(ctx, key1, a, b, c)
		client.SAdd(ctx, key2, c, d, e)

		destination := "dest"

		result, err := client.SInterStore(ctx, destination, key1, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))

		strings, err := client.SMembers(ctx, destination).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{c}))
	})

	It("SMOVE", func() {
		key1, a, b, c, d, e := "SMOVE_key1", "a", "b", "c", "d", "e"
		key2 := "SMOVE_key2"
		client.SAdd(ctx, key1, a, b, c, d)
		client.SAdd(ctx, key2, c, d, e)

		result, err := client.SMove(ctx, key1, key2, a).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(true))

		result, err = client.SIsMember(ctx, key1, a).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(false))

		result, err = client.SIsMember(ctx, key2, a).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(true))
	})

	It("SPOP", func() {
		key1, a, b, c, d, e := "SPOP_key1", "a", "b", "c", "d", "e"
		client.SAdd(ctx, key1, a, b, c, d, e)

		result, err := client.SPop(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Or(Equal(a), Equal(b), Equal(c), Equal(d), Equal(e)))

		count, err := client.SCard(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(4)))
	})

	It("SRANDMEMBER", func() {
		key1, a, b, c, d, e := "SRANDMEMBER_key1", "a", "b", "c", "d", "e"
		client.SAdd(ctx, key1, a, b, c, d, e)

		result, err := client.SRandMember(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Or(Equal(a), Equal(b), Equal(c), Equal(d), Equal(e)))

		count, err := client.SCard(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(5)))
	})

	It("SUNION", func() {
		key1, a, b, c := "SUNION_key1", "a", "b", "c"
		key2, d, e := "SUNION_key1", "d", "e"
		client.SAdd(ctx, key1, a, b, c)
		client.SAdd(ctx, key2, d, e)

		result, err := client.SUnion(ctx, key1, key2).Result()
		Expect(err).To(BeNil())

		expected := []string{a, b, c, d, e}
		sort.Strings(result)
		sort.Strings(expected)
		Expect(result).To(BeEquivalentTo(expected))
	})

	It("SUNIONSTORE", func() {
		key1, a, b, c, d, e := "SUNIONSTORE_key1", "a", "b", "c", "d", "e"
		key2, f, g, h := "SUNIONSTORE_key2", "f", "g", "h"
		client.SAdd(ctx, key1, a, b, c, d, e)
		client.SAdd(ctx, key2, f, g, h)

		dest := "SUNIONSTORE_dest"

		result, err := client.SUnionStore(ctx, dest, key1, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(8)))

		count, err := client.SCard(ctx, dest).Result()
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(8)))
	})

	PIt("SSCAN", func() {

	})
})
