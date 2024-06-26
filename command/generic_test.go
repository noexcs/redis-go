package command_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math"
	"sort"
	"time"
)

var _ = Describe("generic", func() {
	It("keys", func() {
		client.FlushDB(ctx)

		client.MSet(ctx, "firstname", "John", "lastname", "Doe", "age", "21")
		result, err := client.Keys(ctx, "*name*").Result()
		Expect(err).To(BeNil())
		sort.Strings(result)
		Expect(result).To(Equal([]string{"firstname", "lastname"}))

		strings, err := client.Keys(ctx, "a??").Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{"age"}))

		result, err = client.Keys(ctx, "*").Result()
		Expect(err).To(BeNil())
		sort.Strings(result)
		Expect(result).To(Equal([]string{"age", "firstname", "lastname"}))
	})

	key := "generic_key"
	v1 := "v1"

	It("PING", func() {
		pong, err := client.Ping(ctx).Result()
		Expect(err).To(BeNil())
		Expect(pong).To(Equal("PONG"))
	})

	It("FLUSHDB", func() {
		err := client.Set(ctx, key, v1, 0).Err()
		Expect(err).To(BeNil())

		client.FlushDB(ctx)

		val, err := client.Get(ctx, key).Result()
		Expect(err).NotTo(BeNil())
		Expect(val).To(Equal(""))
	})

	It("copy", func() {
		key := "generic_copy_key"
		client.Set(ctx, key, "sheep", 0)

		destKey := "generic_copy_key_copy"
		client.Copy(ctx, key, destKey, 0, false)
		result, err := client.Get(ctx, destKey).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("sheep"))
	})

	It("del", func() {
		key1 := "generic_del_key1"
		key2 := "generic_del_key2"
		key3 := "generic_del_key3"
		client.Set(ctx, key1, "sheep", 0)
		client.Set(ctx, key2, "sheep", 0)

		count, err := client.Del(ctx, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(count).To(Equal(int64(2)))

		result, err := client.MGet(ctx, key1, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal([]any{nil, nil}))
	})

	It("exists", func() {
		key := "generic_exists_key"
		client.Set(ctx, key, "sheep", 0)

		result, err := client.Exists(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))
	})

	It("dump", func() {
		key := "generic_dump_key"
		client.Set(ctx, key, 10, 0)

		result, err := client.Dump(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("\x00\xc0\n\n\x00n\x9fWE\x0e\xaec\xbb"))
	})

	It("exists", func() {
		key1 := "generic_expire_key1"
		key2 := "generic_expire_key2"
		key3 := "generic_expire_key3"

		client.Set(ctx, key1, "sheep", 0)
		client.Set(ctx, key2, "sheep", 0)

		result, err := client.Exists(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(1)))

		result, err = client.Exists(ctx, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(0)))

		result, err = client.Exists(ctx, key1, key2, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(int64(2)))
	})

	It("RENAME", func() {
		key := "rename_key"
		newKey := "rename_key_new"
		value := "Hello RENAME"
		client.Set(ctx, key, value, 0)
		client.Rename(ctx, key, newKey)

		result, err := client.Get(ctx, newKey).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal(value))
	})

	It("RENAMENX", func() {
		key := "rename_key"
		newKey := "rename_key_new"
		value := "Hello RENAME"

		client.Set(ctx, key, value, 0)
		client.Set(ctx, newKey, "World", 0)

		client.RenameNX(ctx, key, newKey)

		result, err := client.Get(ctx, newKey).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("World"))
	})

	It("RESTORE", func() {
		key := "restore_key"
		client.Del(ctx, key)
		serialization := "\n\x17\x17\x00\x00\x00\x12\x00\x00\x00\x03\x00\x00\xc0\x01\x00\x04\xc0\x02\x00\x04\xc0\x03\x00\xff\x04\x00u#<\xc0;.\xe9\xdd"
		s, err2 := client.Restore(ctx, key, 0, serialization).Result()
		Expect(err2).To(BeNil())
		Expect(s).To(Equal("OK"))

		result, err := client.Type(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("list"))

		strings, err := client.LRange(ctx, key, 0, -1).Result()
		Expect(err).To(BeNil())
		Expect(strings).To(Equal([]string{"1", "2", "3"}))
	})

	It("TYPE", func() {
		key1 := "TYPE_key_string"
		key2 := "TYPE_key_list"
		key3 := "TYPE_key_set"

		client.Set(ctx, key1, "value", 0)
		client.LPush(ctx, key2, "value")
		client.SAdd(ctx, key3, "value")

		result, err := client.Type(ctx, key1).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("string"))

		result, err = client.Type(ctx, key2).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("list"))

		result, err = client.Type(ctx, key3).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("set"))
	})

	It("EXPIRE", func() {
		key := "expire_key"

		client.Set(ctx, key, "value", 0)
		client.Expire(ctx, key, 5*time.Second)

		result, err := client.Get(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).To(Equal("value"))

		time.Sleep(5 * time.Second)
		s, err := client.Get(ctx, key).Result()
		Expect(err).NotTo(BeNil())
		Expect(s).To(Equal(""))
	})

	It("TTL", func() {
		key := "ttl_key"
		client.Set(ctx, key, "value", 5*time.Second)

		time.Sleep(3 * time.Second)

		result, err := client.TTL(ctx, key).Result()
		Expect(err).To(BeNil())
		Expect(result).Should(Satisfy(func(t time.Duration) bool {
			if math.Abs(t.Seconds()-2) < 0.01 {
				return true
			}
			return false
		}))

	})
})
