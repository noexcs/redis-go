package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/datastruct"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

// 命令			描述

// SADD			向集合添加一个或多个成员
// SCARD			获取集合的成员数
// SDIFF			返回给定所有集合的差集
// SDIFFSTORE	返回给定所有集合的差集并存储在 destination 中
// SINTER		返回给定所有集合的交集
// SINTERSTORE	返回给定所有集合的交集并存储在 destination 中
// SISMEMBER		判断 member 元素是否是集合 key 的成员
// SMEMBERS		返回集合中的所有成员
// SMOVE			将 member 元素从 source 集合移动到 destination 集合
// SPOP			移除并返回集合中的一个随机元素
// SRANDMEMBER	返回集合中一个或多个随机数
// SREM			移除集合中一个或多个成员
// SUNION		返回所有给定集合的并集
// SUNIONSTORE	所有给定集合的并集存储在 destination 集合中
// SSCAN			迭代集合中的元素

func init() {
	RegisterCommand("sadd", execSadd, nil, nil, -2, FlagReadonly)
	RegisterCommand("srem", execSrem, nil, nil, -3, FlagReadonly)
	RegisterCommand("smembers", execSmember, nil, nil, 2, FlagReadonly)

	RegisterCommand("sismember", execSismember, nil, nil, 3, FlagReadonly)
	RegisterCommand("scard", execScard, nil, nil, 2, FlagReadonly)
}

// Redis SISMEMBER 用于判断元素 member 是否集合 key 的成员。
// SISMEMBER KEY MEMBER
// 返回值 整数:
// 1 如果成员元素是集合的成员，返回 1 。
// 0 如果成员元素不是集合的成员，或 key 不存在，返回 0 。
func execSismember(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	set, errResponse := getOrInitSet(db, key)
	if errResponse != nil {
		return errResponse
	}
	code := 0
	if set.Contains((*args.Data[2]).String()) {
		code = 1
	}
	return &parser.Response{Args: &resp2.Integer{Data: int64(code)}}
}

// SMEMBERS key
// Return
// Array reply: all elements of the set.
func execSmember(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	set, errResponse := getOrInitSet(db, key)
	if errResponse != nil {
		return errResponse
	}
	members := set.Members()
	memberCnt := len(members)

	result := resp2.Array{
		Data:   make([]*resp2.RespType, memberCnt),
		Length: memberCnt,
	}

	for i := 0; i < memberCnt; i++ {
		var r resp2.RespType = &resp2.BulkString{Data: []byte(members[i])}
		result.Data[i] = &r
	}

	return &parser.Response{Args: &result}
}

// Redis SCARD 命令返回集合中元素的数量。
// SCARD KEY_NAME
// 返回值 整数: 集合中成员的数量。
// 当集合 key 不存在时，返回 0 。
func execScard(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	set, errResponse := getOrInitSet(db, key)
	if errResponse != nil {
		return errResponse
	}
	size := set.Size()
	return &parser.Response{Args: &resp2.Integer{Data: int64(size)}}
}

// SREM key member [member ...]
// Return
// Integer reply: the number of members that were removed from the set,
//
//	not including non existing members.
func execSrem(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	set, ErrResponse := getOrInitSet(db, key)
	if ErrResponse != nil {
		return ErrResponse
	}

	count := 0
	for i := 2; i < args.Length; i++ {
		v := (*args.Data[i]).String()
		if set.Contains(v) {
			set.Remove(v)
			count++
		}
	}

	return &parser.Response{Args: &resp2.Integer{Data: int64(count)}}
}

// SADD key member [member ...]
// Return
// Integer reply: the number of elements that were added to the set,
//
//	not including all the elements already present in the set.
func execSadd(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	set, errResponse := getOrInitSet(db, key)
	if errResponse != nil {
		return errResponse
	}
	var count int64
	for i := 2; i < args.Length; i++ {
		set.Add((*args.Data[i]).String())
		count++
	}

	return &parser.Response{Args: &resp2.Integer{Data: count}}
}

// if the value is not the type set, return WRONGTYPE Error
func getOrInitSet(db database.DB, key string) (*datastruct.Set, *parser.Response) {
	value, exist := db.GetValue(key)
	if !exist {
		set := datastruct.NewSet()
		db.SetValue(key, set)
		return set, nil
	}

	set, ok := value.(*datastruct.Set)
	if !ok {
		return nil, &parser.Response{Args: nil, Err: &parser.Error{
			Kind:    "WRONGTYPE",
			Message: "Operation against a key holding the wrong kind of value",
		}}
	}
	return set, nil
}
