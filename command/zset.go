package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/datastruct"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

// 命令				描述
// ZADD				向有序集合添加一个或多个成员，或者更新已存在成员的分数
// ZCARD				获取有序集合的成员数
// ZCOUNT			计算在有序集合中指定区间分数的成员数
// ZINCRBY			有序集合中对指定成员的分数加上增量 increment
// ZINTERSTORE		计算给定的一个或多个有序集的交集并将结果集存储在新的有序集合 key 中
// ZLEXCOUNT			在有序集合中计算指定字典区间内成员数量
// ZRANGE			通过索引区间返回有序集合成指定区间内的成员
// ZRANGEBYLEX		通过字典区间返回有序集合的成员
// ZRANGEBYSCORE		通过分数返回有序集合指定区间内的成员
// ZRANK				返回有序集合中指定成员的索引
// ZREM				移除有序集合中的一个或多个成员
// ZREMRANGEBYLEX	移除有序集合中给定的字典区间的所有成员
// ZREMRANGEBYRANK	移除有序集合中给定的排名区间的所有成员
// ZREMRANGEBYSCORE	移除有序集合中给定的分数区间的所有成员
// ZREVRANGE			返回有序集中指定区间内的成员，通过索引，分数从高到底
// ZREVRANGEBYSCORE	返回有序集中指定分数区间内的成员，分数从高到低排序
// ZREVRANK			返回有序集合中指定成员的排名，有序集成员按分数值递减(从大到小)排序
// ZSCORE			返回有序集中，成员的分数值
// ZUNIONSTORE		计算一个或多个有序集的并集，并存储在新的 key 中
// ZSCAN				迭代有序集合中的元素（包括元素成员和元素分值）

func init() {
	RegisterCommand("ZADD", execZadd, nil, nil, -4, FlagWrite)
}

// execZadd ZADD key score member [score member...]
func execZadd(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	sortedset, response, _ := getOrInitSortedset(db, key, true)
	if response != nil {
		return response
	}
	count := int64(0)
	for i := 2; i+1 < len(args.Data); i += 2 {
		score := *args.Data[i]
		//todo: validate the score type
		if s, ok := score.(*resp2.Integer); ok {
			sortedset.Add(s.Data, (*args.Data[i+1]).String())
		}
		count++
	}

	return &parser.Response{Args: &resp2.Integer{Data: count}}
}

func getOrInitSortedset(db database.DB, key string, init bool) (*datastruct.Sortedset, *parser.Response, bool) {
	value, exist := db.GetValue(key)
	if !exist {
		if init {
			newSortedset := datastruct.NewSortedset()
			db.SetValue(key, newSortedset)
			return newSortedset, nil, false
		} else {
			return nil, nil, false
		}
	}

	sortedset, ok := value.(*datastruct.Sortedset)
	if !ok {
		return nil, &parser.Response{Args: nil, Err: &parser.Error{
			Kind:    "WRONGTYPE",
			Message: "Operation against a key holding the wrong kind of value",
		}}, true
	}
	return sortedset, nil, true
}
