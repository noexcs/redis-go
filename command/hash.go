package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/datastruct"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

//命令	    说明

//HDEL		用于删除哈希表中一个或多个字段
//HEXISTS	用于判断哈希表中字段是否存在
//HGET		获取存储在哈希表中指定字段的值
//HGETALL	获取在哈希表中指定 key 的所有字段和值
//HINCRBY	为存储在 key 中的哈希表指定字段做整数增量运算
//HKEYS		获取存储在 key 中的哈希表的所有字段
//HLEN		获取存储在 key 中的哈希表的字段数量
//HSET		用于设置存储在 key 中的哈希表字段的值
//HVALS		用于获取哈希表中的所有值

func init() {
	RegisterCommand("hset", execHset, nil, nil, -4, FlagReadonly)
	RegisterCommand("hget", execHget, nil, nil, 3, FlagReadonly)
}

// Redis HGET 命令用于返回哈希表中指定字段 field 的值。
// HGET key field
// 返回值
// 多行字符串: 返回给定字段的值。如果给定的字段或 key 不存在时，返回 nil 。
func execHget(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	hashmap, errResponse, exist := getOrInitHashmap(db, key, false)
	if !exist {
		return &parser.Response{Args: resp2.MakeNullBulkString()}
	}
	if errResponse != nil {
		return errResponse
	}
	field := (*args.Data[2]).String()
	v, found := hashmap.Get(field)
	if !found {
		return &parser.Response{Args: resp2.MakeNullBulkString()}
	}
	return &parser.Response{Args: &resp2.SimpleString{Data: v}}
}

// Redis Hset 命令用于为存储在 key 中的哈希表的 field 字段赋值 value 。
// 如果哈希表不存在，一个新的哈希表被创建并进行 HSET 操作。
// 如果字段（field）已经存在于哈希表中，旧值将被覆盖。
// HSET key field value [field value ...]
// 返回值 整数:
// 只有被修改返回0 ，有增加返回增加的 field 个数。
func execHset(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	hashmap, errResponse, _ := getOrInitHashmap(db, key, true)
	if errResponse != nil {
		return errResponse
	}

	var count int64 = 0
	for i := 2; i+1 < args.Length; i += 2 {
		field := (*args.Data[i]).String()
		value := (*args.Data[i+1]).String()
		if !hashmap.Contains(field) {
			count++
		}
		hashmap.Put(field, value)
	}

	return &parser.Response{Args: &resp2.Integer{Data: count}}
}

// if the value is not the type set, return WRONGTYPE Error
func getOrInitHashmap(db database.DB, key string, init bool) (*datastruct.Hashmap, *parser.Response, bool) {
	value, exist := db.GetValue(key)
	if !exist {
		if init {
			newHashmap := datastruct.NewHashmap()
			db.SetValue(key, newHashmap)
			return newHashmap, nil, false
		} else {
			return nil, nil, false
		}
	}

	hashmap, ok := value.(*datastruct.Hashmap)
	if !ok {
		return nil, &parser.Response{Args: nil, Err: &parser.Error{
			Kind:    "WRONGTYPE",
			Message: "Operation against a key holding the wrong kind of value",
		}}, true
	}
	return hashmap, nil, true
}
