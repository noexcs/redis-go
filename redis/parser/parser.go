package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"io"
	"net"
	"strconv"
)

// ParseIncomeStream
// 从连接中读取客户端请求，并发送到requestChan channel中，
// 直到断开连接，关闭管道
func ParseIncomeStream(conn net.Conn) <-chan *Request {
	requestChan := make(chan *Request)
	go func() {
		defer close(requestChan)
		err := parseBackground(conn, requestChan)
		if err != nil {
			log.Debug("Client message parse error: ", err)
		}
	}()
	return requestChan
}

// 解析 Redis 数据
// 参考：https://redis.io/docs/latest/develop/reference/protocol-spec
func parseBackground(conn net.Conn, requestChan chan *Request) error {
	reader := bufio.NewReader(conn)
	for {
		// 从连接中读取数据，直到遇到定界符（'\n'）
		// 返回一个slice，包含数据以及定界符（'\n'）
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// If ReadBytes encounters an error before finding a delimiter,
			// it returns the data read before the error and the error itself (often io.EOF).
			requestChan <- MakeErrorRequest("ERR", "ReadBytes error.")
			return err
		}
		length := len(line)
		if length < 3 || line[length-2] != '\r' {
			requestChan <- MakeErrorRequest("ERR", "Illegal resp2 type line.")
			continue
		}

		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+':
			// line: +OK
			requestChan <- parseSimpleString(line)
		case '-':
			// line: -ERR message
			requestChan <- parseSimpleError(line)
		case ':':
			// line: :1000
			requestChan <- parseInteger(line)
		case '$':
			// line: $5
			requestChan <- parseBulkString(line, reader)
		case '*':
			// line: *3
			requestChan <- parseArray(line, reader, requestChan)
		default:
			requestChan <- MakeErrorRequest("ERR", "Encounter error while parsing header bytes.")
		}
	}
}

func parseInteger(line []byte) *Request {
	integer, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil {
		return MakeErrorRequest("ERR", "Illegal Integer resp2 type encoding.")
	}
	return &Request{
		Args: &resp2.Integer{Data: integer},
	}
}

func parseSimpleError(line []byte) *Request {
	index := bytes.Index(line, []byte(" "))
	kind := line[1:index]
	data := line[index+1:]
	return &Request{
		Args: &resp2.SimpleError{Kind: string(kind), Data: string(data)},
	}
}

func parseSimpleString(line []byte) *Request {
	content := string(line[1:])
	return &Request{
		Args: &resp2.SimpleString{Data: content},
	}
}

func parseBulkString(line []byte, reader io.Reader) *Request {
	strLen, err := strconv.ParseInt(string(line[1:]), 10, 64)
	if err != nil || strLen < -1 {
		return MakeErrorRequest("ERR", "Illegal bulk string header.")
	} else if strLen == -1 {
		return &Request{
			Args: &resp2.BulkString{Data: nil},
		}
	} else {
		// +2 means '\r\n'
		body := make([]byte, strLen+2)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return MakeErrorRequest("ERR", "Encounter error while reading bytes from bulk string.")
		}
		return &Request{
			// remove '\r\n'
			Args: &resp2.BulkString{Data: body[:strLen]},
		}
	}
}

func parseArray(header []byte, reader *bufio.Reader, requestChan chan *Request) *Request {

	// header: *length
	if len(header) < 2 {
		return MakeErrorRequest("ERR", "Encounter error while parsing array header.")
	}
	length, err := strconv.ParseInt(string(header[1:]), 10, 32)
	if err != nil {
		return MakeErrorRequest("ERR", "Encounter error while parsing array length.")
	}
	arr := &resp2.Array{}

	arr.Length = int(length)
	arr.Data = make([]*resp2.RespType, length)

	for i := 0; i < arr.Length; i++ {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return MakeErrorRequest("ERR", "Encounter error while reading bytes from array.")
		}

		length := len(line)
		if length < 3 || line[length-2] != '\r' {
			return MakeErrorRequest("ERR", "Illegal resp2 type line.")
		}

		var currentItem *resp2.RespType
		var parseResult *Request
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		switch line[0] {
		case '+':
			// line: +OK
			parseResult = parseSimpleString(line)
		case '-':
			// line: -ERR message
			parseResult = parseSimpleError(line)
		case ':':
			// line: :1000
			parseResult = parseInteger(line)
		case '$':
			// line: $5
			parseResult = parseBulkString(line, reader)
		case '*':
			// line: *3
			parseResult = parseArray(line, reader, requestChan)
		default:
			parseResult = MakeErrorRequest("ERR", "Encounter error while parsing header bytes.")
		}

		if parseResult.Err != nil {
			return parseResult
		}
		currentItem = &parseResult.Args

		arr.Data[i] = currentItem
	}

	return &Request{
		Args: arr,
		Err:  nil,
	}
}

type Error struct {
	Kind    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error kind: %s, Error data: %s", e.Kind, e.Message)
}
