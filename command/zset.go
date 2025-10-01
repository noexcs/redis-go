package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/datastruct"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

func init() {
	// ZADD 向有序集合添加一个或多个成员，或者更新已存在成员的分数
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
