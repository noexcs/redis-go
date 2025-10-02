package resp2

import (
	"bufio"
	"bytes"
	"io"
	"strconv"

	"github.com/noexcs/redis-go/redis/parser/resp"
)

type Resp2Parser struct{}

func (p *Resp2Parser) Parse(reader *bufio.Reader) (resp.RespValue, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	length := len(line)
	if length < 3 || line[length-2] != '\r' {
		return nil, ErrProtocol
	}
	line = bytes.TrimSuffix(line, []byte{'\r', '\n'})

	switch line[0] {
	case '+':
		return parseSimpleString(line), nil
	case '-':
		return parseSimpleError(line), nil
	case ':':
		return parseInteger(line), nil
	case '$':
		return parseBulkString(line, reader), nil
	case '*':
		return parseArray(line, reader), nil
	default:
		return nil, ErrProtocol
	}
}

var ErrProtocol = io.ErrUnexpectedEOF

func parseInteger(line []byte) resp.RespValue {
	integer, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil {
		return &SimpleError{Kind: "ERR", Data: "Illegal Integer"}
	}
	return &Integer{Data: integer}
}

func parseSimpleError(line []byte) resp.RespValue {
	index := bytes.Index(line, []byte(" "))
	if index < 0 {
		return &SimpleError{Kind: "ERR", Data: "Malformed error"}
	}
	kind := line[1:index]
	data := line[index+1:]
	return &SimpleError{Kind: string(kind), Data: string(data)}
}

func parseSimpleString(line []byte) resp.RespValue {
	return &SimpleString{Data: string(line[1:])}
}

func parseBulkString(line []byte, reader io.Reader) resp.RespValue {
	strLen, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil || strLen < -1 {
		return &SimpleError{Kind: "ERR", Data: "Illegal bulk string header"}
	} else if strLen == -1 {
		return &BulkString{Data: nil}
	}
	body := make([]byte, strLen+2)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return &SimpleError{Kind: "ERR", Data: "Error reading bulk string"}
	}
	return &BulkString{Data: body[:strLen]}
}

func parseArray(header []byte, reader *bufio.Reader) resp.RespValue {
	length, err := strconv.Atoi(string(header[1:]))
	if err != nil {
		return &SimpleError{Kind: "ERR", Data: "Invalid array length"}
	}
	arr := &Array{
		Length: length,
		Data:   make([]resp.RespValue, length),
	}
	for i := 0; i < length; i++ {
		val, err := (&Resp2Parser{}).Parse(reader)
		if err != nil {
			return &SimpleError{Kind: "ERR", Data: "Array parse error"}
		}
		arr.Data[i] = val
	}
	return arr
}
