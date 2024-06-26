package command_test

import (
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("zset", func() {
	Context("zset", func() {

		It("ZADD", func() {
			key := "ZADD_example_key"
			member := &redis.Z{
				Score:  1,
				Member: "one",
			}
			result, err := client.ZAdd(ctx, key, member).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))
		})

		It("ZREM", func() {
			key1 := "ZREM_key1"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2, member3)

			result, err := client.ZRem(ctx, key1, member2.Member.(string)).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))

			strings, err := client.ZRange(ctx, key1, 0, -1).Result()
			Expect(err).To(BeNil())
			expected := []string{"one", "three"}
			sort.Strings(expected)
			sort.Strings(strings)
			Expect(strings).To(Equal(expected))
		})

		It("ZSCORE", func() {
			key := "ZSCORE_key"
			client.ZAdd(ctx, key, &redis.Z{
				Score:  1,
				Member: "one"})
			score, err := client.ZScore(ctx, key, "one").Result()
			Expect(err).To(BeNil())
			Expect(score).To(Equal(float64(1)))
		})

		It("ZUNION", func() {
			key1 := "ZUNION_key1"
			key2 := "ZUNION_key2"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2)
			client.ZAdd(ctx, key2, member1, member2, member3)
			store := redis.ZStore{
				Keys:      []string{key1, key2},
				Weights:   nil,
				Aggregate: "",
			}

			result, err := client.ZUnion(ctx, store).Result()
			Expect(err).To(BeNil())
			expected := []string{member1.Member.(string), member3.Member.(string), member2.Member.(string)}
			sort.Strings(result)
			sort.Strings(expected)
			Expect(result).To(Equal(expected))

			zs, err := client.ZUnionWithScores(ctx, store).Result()
			member1.Score = 2
			member2.Score = 4
			member3.Score = 3
			expectedZs := []redis.Z{*member1, *member3, *member2}
			sort.Slice(zs, func(i, j int) bool {
				return zs[i].Member.(string) < zs[j].Member.(string)
			})
			sort.Slice(expectedZs, func(i, j int) bool {
				return zs[i].Member.(string) < zs[j].Member.(string)
			})
			Expect(err).To(BeNil())
			Expect(zs).To(Equal(expectedZs))

		})

		It("ZUNIONSTORE", func() {
			key1 := "ZUNIONSTORE_key1"
			key2 := "ZUNIONSTORE_key2"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2)
			client.ZAdd(ctx, key2, member1, member2, member3)

			key_out := "ZUNIONSTORE_out"
			store := redis.ZStore{
				Keys:      []string{key1, key2},
				Weights:   []float64{2, 3},
				Aggregate: "",
			}
			client.ZUnionStore(ctx, key_out, &store)

			zs, err := client.ZRangeWithScores(ctx, key_out, 0, -1).Result()
			member1.Score = 5
			member2.Score = 10
			member3.Score = 9
			expectedZs := []redis.Z{*member1, *member3, *member2}
			Expect(err).To(BeNil())
			Expect(zs).To(Equal(expectedZs))
		})

		It("ZRANK", func() {
			key := "ZRANK_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key, member1, member2, member3)

			result, err := client.ZRank(ctx, key, member3.Member.(string)).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(2)))

			result, err = client.ZRank(ctx, key, "four").Result()
			Expect(err).NotTo(BeNil())
			Expect(result).To(Equal(int64(0)))
		})

		It("ZRANDMEMBER", func() {
			key := "ZRANDMEMBER_key"
			members := make([]*redis.Z, 7)
			for i := 0; i < 7; i++ {
				members[i] = &redis.Z{
					Score:  0,
					Member: string(rune('a' + i)),
				}
			}
			client.ZAdd(ctx, key, members...)

			for i := 0; i < 5; i++ {
				result, err := client.ZRandMember(ctx, key, 1, false).Result()
				Expect(err).To(BeNil())
				Expect(result).Should(Satisfy(func(x []string) bool {
					for _, member := range x {
						if !(len(member) == 1 && rune(member[0]) >= 'a' && rune(member[0]) <= 'a'+6) {
							return false
						}
					}
					return true
				}))
			}
		})

		XIt("ZRANGE", func() {
			key := "ZRANGE_key"
			members := make([]*redis.Z, 5)
			for i := 0; i < len(members); i++ {
				members[i] = &redis.Z{
					Score:  0,
					Member: string(rune('a' + i)),
				}
			}
			client.ZAdd(ctx, key, members...)

			result, err := client.ZRange(ctx, key, 0, -1).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]string{"a", "b", "c", "d", "e"}))
		})

		// deprecated
		PIt("ZRANGEBYSCORE", func() {

		})

		It("ZRANGESTORE", func() {
			key := "ZRANGESTORE_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			member4 := &redis.Z{
				Score:  4,
				Member: "four",
			}
			client.ZAdd(ctx, key, member1, member2, member3, member4)

			keyDest := "ZRANGESTORE_dest"
			zRangeArgs := redis.ZRangeArgs{
				Key:   key,
				Start: 2,
				Stop:  -1,
			}
			result, err := client.ZRangeStore(ctx, keyDest, zRangeArgs).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(2)))
		})

		It("ZREMRANGEBYLEX", func() {
			key := "ZREMRANGEBYLEX_key"

			members := make([]*redis.Z, 5)
			for i := 0; i < len(members); i++ {
				members[i] = &redis.Z{
					Score: 0,
				}
			}
			members[0].Member = "aaaa"
			members[1].Member = "b"
			members[2].Member = "c"
			members[3].Member = "d"
			members[4].Member = "e"

			client.ZAdd(ctx, key, members...)

			members1 := make([]*redis.Z, 5)
			for i := 0; i < len(members1); i++ {
				members1[i] = &redis.Z{
					Score: 0,
				}
			}
			members1[0].Member = "foo"
			members1[1].Member = "zap"
			members1[2].Member = "zip"
			members1[3].Member = "ALPHA"
			members1[4].Member = "alpha"
			client.ZAdd(ctx, key, members1...)

			result, err := client.ZRemRangeByLex(ctx, key, "[alpha", "[omega").Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(6)))

			strings, err := client.ZRange(ctx, key, 0, -1).Result()
			Expect(err).To(BeNil())
			expected := []string{"ALPHA", "aaaa", "zap", "zip"}
			sort.Strings(strings)
			sort.Strings(expected)
			Expect(strings).To(Equal(expected))
		})

		It("ZREMRANGEBYRANK", func() {
			key := "ZREMRANGEBYRANK_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key, member1, member2, member3)
			result, err := client.ZRemRangeByRank(ctx, key, 0, 1).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(2)))

			zs, err := client.ZRangeWithScores(ctx, key, 0, -1).Result()
			Expect(err).To(BeNil())
			Expect(zs).To(Equal([]redis.Z{*member3}))
		})

		It("ZREMRANGEBYSCORE", func() {
			key := "ZREMRANGEBYSCORE_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key, member1, member2, member3)
			result, err := client.ZRemRangeByScore(ctx, key, "-inf", "(2").Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))

			zs, err := client.ZRangeWithScores(ctx, key, 0, -1).Result()
			Expect(err).To(BeNil())
			expected := []redis.Z{*member2, *member3}

			Expect(zs).To(Equal(expected))
		})

		// deprecated
		PIt("ZREVRANGE", func() {})
		// deprecated
		PIt("ZREVRANGEBYLEX", func() {})
		// deprecated
		PIt("ZREVRANGEBYSCORE", func() {})

		It("ZREVRANK", func() {
			key := "ZREVRANK_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key, member1, member2, member3)

			rank1, err := client.ZRevRank(ctx, key, member1.Member.(string)).Result()
			Expect(err).To(BeNil())
			Expect(rank1).To(Equal(int64(2)))

			//redis> ZREVRANK myzset "four"
			//(nil)
			rank2, err := client.ZRevRank(ctx, key, "four").Result()
			Expect(err).NotTo(BeNil())
			Expect(rank2).To(Equal(int64(0)))
		})

		It("ZCARD", func() {
			key := "ZCARD_example_key"
			member := &redis.Z{
				Score:  1,
				Member: "one",
			}
			client.ZAdd(ctx, key, member)

			result, err := client.ZCard(ctx, key).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))
		})

		It("ZCOUNT", func() {
			key := "ZCOUNT_example_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			member4 := &redis.Z{
				Score:  4,
				Member: "four",
			}
			member5 := &redis.Z{
				Score:  5,
				Member: "five",
			}
			client.ZAdd(ctx, key, member1, member2, member3, member4, member5)

			result, err := client.ZCount(ctx, key, "2", "5").Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(4)))
		})

		It("ZDIFF", func() {
			key1 := "ZDIFF_key1"
			key2 := "ZDIFF_key2"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2, member3)
			client.ZAdd(ctx, key2, member1, member2)

			result, err := client.ZDiff(ctx, key1, key2).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]string{"three"}))
		})

		It("ZDIFFSTORE", func() {
			key1 := "ZDIFFSTORE_key1"
			key2 := "ZDIFFSTORE_key2"
			key3 := "ZDIFFSTORE_key3"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2, member3)
			client.ZAdd(ctx, key2, member1, member2)

			result, err := client.ZDiffStore(ctx, key3, key1, key2).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(1)))

			client.ZRange(ctx, key3, 0, -1)
			strings, err := client.ZRange(ctx, key3, 0, -1).Result()
			Expect(err).To(BeNil())
			Expect(strings).To(Equal([]string{"three"}))
		})

		It("ZINCRBY", func() {
			key := "ZINCRBY_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			client.ZAdd(ctx, key, member1, member2)

			result, err := client.ZIncrBy(ctx, key, 2, member1.Member.(string)).Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(3.0))

			score, err := client.ZScore(ctx, key, member1.Member.(string)).Result()
			Expect(err).To(BeNil())
			Expect(score).To(Equal(3.0))
		})

		It("ZINTER", func() {
			key1 := "ZINTER_key1"
			key2 := "ZINTER_key2"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2, member3)
			client.ZAdd(ctx, key2, member1, member2)

			store := &redis.ZStore{
				Keys:      []string{key1, key2},
				Weights:   nil,
				Aggregate: "",
			}
			result, err := client.ZInter(ctx, store).Result()
			Expect(err).To(BeNil())
			expected := []string{member1.Member.(string), member2.Member.(string)}
			sort.Strings(result)
			sort.Strings(expected)
			Expect(result).To(Equal(expected))
		})

		It("ZINTERSTORE", func() {
			key1 := "ZINTERCRAD_key1"
			key2 := "ZINTERCARD_key2"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key1, member1, member2, member3)
			client.ZAdd(ctx, key2, member1, member2)

			store := &redis.ZStore{
				Keys:      []string{key1, key2},
				Weights:   nil,
				Aggregate: "",
			}
			destination := "ZINTERSTORE_key"
			client.ZInterStore(ctx, destination, store)
			result, err := client.ZRange(ctx, destination, 0, -1).Result()
			Expect(err).To(BeNil())
			sort.Strings(result)
			expected := []string{member1.Member.(string), member2.Member.(string)}
			sort.Strings(expected)
			Expect(result).To(Equal(expected))
		})

		It("ZLEXCOUNT", func() {
			key := "ZLEXCOUNT_key"
			members := make([]*redis.Z, 7)
			for i := 0; i < 7; i++ {
				members[i] = &redis.Z{
					Score:  0,
					Member: string(rune('a' + i)),
				}
			}
			client.ZAdd(ctx, key, members...)

			result, err := client.ZLexCount(ctx, key, "[b", "[e").Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(int64(4)))
		})

		// deprecated
		PIt("ZRANGEBYLEX", func() {

		})
		// deprecated
		PIt("ZRANGEBYSCORE", func() {

		})

		It("ZMPOP", func() {
			result, err := client.ZPopMin(ctx, "ZMPOP_noSuchKey", 1).Result()
			Expect(err).To(BeNil())
			Expect(result).To(HaveLen(0))

			key := "ZMPOP_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			member3 := &redis.Z{
				Score:  3,
				Member: "three",
			}
			client.ZAdd(ctx, key, member1, member2, member3)

			zs, err := client.ZPopMin(ctx, key, 1).Result()
			Expect(err).To(BeNil())
			Expect(zs).To(HaveLen(1))
			Expect(&zs[0]).To(Equal(member1))

			i, err := client.ZPopMax(ctx, key, 1).Result()
			Expect(err).To(BeNil())
			Expect(i).Should(HaveLen(1))
			Expect(i[0]).To(Equal(*member3))
		})

		It("ZMSCORE", func() {
			key := "ZMSCORE_key"
			member1 := &redis.Z{
				Score:  1,
				Member: "one",
			}
			member2 := &redis.Z{
				Score:  2,
				Member: "two",
			}
			client.ZAdd(ctx, key, member1, member2)

			result, err := client.ZMScore(ctx, key, member1.Member.(string), member2.Member.(string), "nofield").Result()
			Expect(err).To(BeNil())
			Expect(result).To(Equal([]float64{1, 2, 0}))
		})

		//ZPOPMAX
		//ZPOPMIN

		//ZSCAN
	})
})
