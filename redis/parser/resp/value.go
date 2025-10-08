package resp

import (
	"bufio"
	"fmt"
)

// RespValue 是所有 RESP 协议类型的统一接口
type RespValue interface {
	ToBytes() []byte
	String() string
}

type RespParser interface {
	Parse(reader *bufio.Reader) (RespValue, error)
	MakeSimpleString(data string) RespValue
}

type Error struct {
	Kind    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error kind: %s, Error data: %s", e.Kind, e.Message)
}
