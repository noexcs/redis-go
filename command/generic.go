package command

import (
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

func init() {
	RegisterCommand("ping", execPing, nil, nil, -1, FlagReadonly)
	RegisterCommand("auth", execAuth, nil, nil, -2, FlagReadonly)
	RegisterCommand("flushdb", execFlushDb, nil, nil, 1, FlagWrite)
}

func execFlushDb(db database.DB, args *resp2.Array) *parser.Response {
	db.FlushDb()
	return &parser.Response{Args: resp2.MakeOKSimpleString()}
}

// Original: AUTH [username] password
// Implementation: AUTH password
// https://redis.io/commands/auth/
func execAuth(db database.DB, args *resp2.Array) *parser.Response {
	//config.Properties.
	password := (*args.Data[1]).String()
	if config.Properties.Requirepass != "" {
		if config.Properties.Requirepass != password {
			return &parser.Response{Args: resp2.MakeSimpleError("ERR ", "invalid password")}
		} else {
			return &parser.Response{Args: resp2.MakeOKSimpleString()}
		}
	} else {
		return &parser.Response{Args: resp2.MakeSimpleError("ERR", "Client sent AUTH, but no password is set")}
	}

	return nil
}

// PING [message]
// PING https://redis.io/commands/ping/
func execPing(db database.DB, args *resp2.Array) *parser.Response {

	if args.Length >= 2 {
		data := (*args.Data[1]).(*resp2.BulkString).Data
		r := resp2.BulkString{Data: data}
		return &parser.Response{Args: &r}
	}

	return &parser.Response{Args: resp2.MakePONGSimpleString()}
}
