package command

import (
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"strings"
)

// redis command refer: https://redis.io/docs/reference/command-tips/
// redis command refer: https://redis.io/docs/reference/key-specs/
// redis command refer: https://redis.com.cn/commands.html

var CmdTable = make(map[string]*command)

type command struct {
	Executor ExecFunc
	Prepare  PreFunc // return related keys command
	Undo     UndoFunc

	// arity 为命令的参数数量，
	// 为正数时，参数数量必须和cmd.arity一致
	// 为负数时，参数数量至少为 -cmd.arity
	Arity int // allow number of args, arity < 0 means len(args) >= -arity
	Flags Flag
}

type Flag uint32

const (
	FlagWrite Flag = 1 << iota
	FlagReadonly
)

type ExecFunc func(db database.DB, args *resp2.Array) *parser.Response
type PreFunc func(db database.DB, args *resp2.Array) *parser.Response
type UndoFunc func(db database.DB, args *resp2.Array) *parser.Response

func RegisterCommand(name string, execFunc ExecFunc, preFunc PreFunc, undoFunc UndoFunc, arity int, flags Flag) {
	cmd := &command{
		Executor: execFunc,
		Prepare:  preFunc,
		Undo:     undoFunc,
		Arity:    arity,
		Flags:    flags,
	}
	CmdTable[strings.ToUpper(name)] = cmd
	//log.WithLocation(fmt.Sprintf("Command %s registered.", name))
}
