package command

import (
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

func init() {
	RegisterCommand("hello", execHello, nil, nil, -2, FlagWrite)
	RegisterCommand("ping", execPing, nil, nil, -1, FlagReadonly)
	RegisterCommand("auth", execAuth, nil, nil, -2, FlagReadonly)
	RegisterCommand("flushdb", execFlushDb, nil, nil, 1, FlagWrite)
	RegisterCommand("del", execDel, nil, nil, -2, FlagWrite)
}

func execHello(db database.DB, args []resp.RespValue) *parser.Response {
	version := args[1].String()
	if version != "2" {
		return &parser.Response{
			Args: resp2.MakeSimpleError("NOPROTO", "sorry this protocol version is not supported"),
			Err:  nil,
		}
	}
	return nil
}

func execFlushDb(db database.DB, args []resp.RespValue) *parser.Response {
	db.FlushDb()
	return &parser.Response{Args: resp2.MakeOKSimpleString()}
}

func execDel(db database.DB, args []resp.RespValue) *parser.Response {
	count := 0
	for i := 1; i < len(args); i++ {
		key := args[i].String()
		if db.Delete(key) {
			count++
		}
	}

	// 主动触发一次过期键清理
	if simpleDB, ok := db.(interface{ DeleteExpiredKeys() }); ok {
		simpleDB.DeleteExpiredKeys()
	}

	return &parser.Response{Args: &resp2.Integer{Data: int64(count)}}
}

func execPing(db database.DB, args []resp.RespValue) *parser.Response {
	if len(args) == 1 {
		return &parser.Response{Args: resp2.MakePONGSimpleString()}
	} else {
		message := args[1].String()
		return &parser.Response{Args: resp2.MakeSimpleString(message)}
	}
}

func execAuth(db database.DB, args []resp.RespValue) *parser.Response {
	if len(config.Properties.RequirePass) == 0 {
		return &parser.Response{Args: resp2.MakeSimpleError("ERR", "Client sent AUTH, but no password is set")}
	}
	pass := args[1].String()
	if pass != config.Properties.RequirePass {
		return &parser.Response{Args: resp2.MakeSimpleError("ERR", "invalid password")}
	}
	return &parser.Response{Args: resp2.MakeOKSimpleString()}
}
