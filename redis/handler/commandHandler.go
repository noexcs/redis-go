package handler

import (
	"github.com/noexcs/redis-go/command"
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/client"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strings"
)

// HandleCommand
// A client sends the Redis server an array consisting of only bulk strings.
// A Redis server replies to clients, sending any valid RESP data type as a reply.
// https://redis.io/docs/reference/protocol-spec/#sending-commands-to-a-redis-server
func HandleCommand(client *client.Client, args resp2.RespType, db database.DB) (result resp2.RespType) {
	log.Debug("Client command: ", args.String())

	// 是否为空命令
	array := args.(*resp2.Array)
	if array.Length < 1 {
		return &resp2.SimpleError{Kind: "ERR", Data: "command is empty"}
	}

	// 是否为不存在的命令
	cmdName := strings.ToUpper((*(*array).Data[0]).String())
	cmd := command.CmdTable[cmdName]
	if cmd == nil {
		return &resp2.SimpleError{Kind: "ERR", Data: "command " + cmdName + " not found"}
	}

	// 命令参数是否足够
	if ok, errResponse := ValidateArity(array.Length, cmd.Arity, cmdName); !ok {
		return errResponse
	}

	// 是否需要密码并已验证
	if config.Properties.RequirePass != "" && !client.Authenticated {
		if cmdName != "AUTH" {
			return &resp2.SimpleError{Kind: "Err", Data: "NOAUTH"}
		}
	}

	response := cmd.Executor(db, array)

	if response != nil {
		if response.Err != nil {
			return &resp2.SimpleError{
				Kind: response.Err.Kind,
				Data: response.Err.Message,
			}
		}
		if response.Args != nil {
			if cmdName == "AUTH" {
				simpleString, ok := response.Args.(*resp2.SimpleString)
				if ok && simpleString.Data == "OK" {
					client.Authenticated = true
				}
			}
			return response.Args
		}
	}

	return &resp2.SimpleString{Data: "OK"}
}

func ValidateArity(argsLen int, arity int, cmdName string) (ok bool, simpleError *resp2.SimpleError) {
	if arity > 0 {
		if argsLen != arity {
			return false, &resp2.SimpleError{Kind: "ERR", Data: "WRONG number of arguments for the '" + cmdName + "' command."}
		}

	}
	if arity < 0 {
		if argsLen < -arity {
			return false, &resp2.SimpleError{Kind: "ERR", Data: "WRONG number of arguments for the '" + cmdName + "' command."}
		}
	}
	return true, nil
}
