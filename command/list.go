package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/datastruct"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strconv"
)

//命令			描述

func init() {
	RegisterCommand("lpush", execLpush, nil, nil, -3, FlagWrite)
	RegisterCommand("rpush", execRpush, nil, nil, -3, FlagWrite)

	RegisterCommand("lpop", execLpop, nil, nil, -2, FlagWrite)
	RegisterCommand("rpop", execRpop, nil, nil, -2, FlagWrite)
}

func execRpop(db database.DB, args *resp2.Array) *parser.Response {
	return pop(db, args, true)
}

// LPOP key [count]
// Available since:
// 1.0.0
// Time complexity:
// O(N) where N is the number of elements returned
// ACL categories:
// @write, @list, @fast
// Removes and returns the first elements of the list stored at key.
//
// By default, the command pops a single element from the beginning of the list.
// When provided with the optional count argument, the reply will consist of up to count elements, depending on the list's length.
//
// Return
// When called without the count argument:
//
// Bulk string reply: the value of the first element, or nil when key does not exist.
//
// When called with the count argument:
//
// Array reply: list of popped elements, or nil when key does not exist.
func execLpop(db database.DB, args *resp2.Array) *parser.Response {
	return pop(db, args, false)
}

func pop(db database.DB, args *resp2.Array, right bool) *parser.Response {
	key := (*args.Data[1]).String()
	list, errResponse, keyExist := getOrInitList(db, key, false)
	if !keyExist {
		return &parser.Response{Args: resp2.MakeNullBulkString()}
	}
	if errResponse != nil {
		return errResponse
	}
	if args.Length >= 3 {
		count, err := strconv.Atoi((*args.Data[2]).String())
		if err != nil {
			return &parser.Response{Args: resp2.MakeSimpleError("ERR", "Wrong count arg.")}
		}
		if int(list.Size()) < count {
			count = int(list.Size())
		}

		arr := resp2.Array{
			Data:   make([]*resp2.RespType, count),
			Length: count,
		}
		for i := 0; i < count; i++ {
			var v string
			var exist bool
			if right {
				v, exist = list.PopRight()
			} else {
				v, exist = list.PopLeft()
			}
			if exist {
				var r resp2.RespType = &resp2.BulkString{Data: []byte(v)}
				arr.Data[i] = &r
			}
		}
		return &parser.Response{Args: &arr}
	} else {
		var first string
		var exist bool
		if right {
			first, exist = list.PopRight()
		} else {
			first, exist = list.PopLeft()
		}
		if !exist {
			return &parser.Response{Args: resp2.MakeNullBulkString()}
		}
		return &parser.Response{Args: &resp2.BulkString{Data: []byte(first)}}
	}
}

// LPUSH key element [element ...]
// Return
// Integer reply: the length of the list after the push operations.
func execLpush(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	list, errResponse, _ := getOrInitList(db, key, true)
	if errResponse != nil {
		return errResponse
	}

	for i := 2; i < args.Length; i++ {
		element := (*args.Data[i]).String()
		list.PushLeft(element)
	}
	return &parser.Response{Args: &resp2.Integer{Data: list.Size()}}
}

// RPUSH key element [element ...]
// Return
// Integer reply: the length of the list after the push operation.
func execRpush(db database.DB, args *resp2.Array) *parser.Response {
	key := (*args.Data[1]).String()
	list, errResponse, _ := getOrInitList(db, key, true)
	if errResponse != nil {
		return errResponse
	}

	for i := 2; i < args.Length; i++ {
		element := (*args.Data[i]).String()
		list.PushRight(element)
	}
	return &parser.Response{Args: &resp2.Integer{Data: list.Size()}}
}

func getOrInitList(db database.DB, key string, init bool) (*datastruct.List, *parser.Response, bool) {
	value, exist := db.GetValue(key)
	if !exist {
		if init {
			newList := datastruct.NewList()
			db.SetValue(key, newList)
			return newList, nil, false
		} else {
			return nil, nil, false
		}
	}

	list, ok := value.(*datastruct.List)
	if !ok {
		return nil, &parser.Response{Args: nil, Err: &parser.Error{
			Kind:    "WRONGTYPE",
			Message: "Operation against a key holding the wrong kind of value",
		}}, true
	}
	return list, nil, true
}
