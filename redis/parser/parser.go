package parser

import (
	"bufio"
	"net"

	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
)

// RespParser 接口
type RespParser interface {
	Parse(reader *bufio.Reader) (resp.RespValue, error)
}

func ParseIncomeStream(conn net.Conn) <-chan *Request {
	requestChan := make(chan *Request)
	go func() {
		defer close(requestChan)
		reader := bufio.NewReader(conn)

		parser := &resp2.Resp2Parser{} // 暂时用 RESP2，将来可以换 RESP3

		for {
			val, err := parser.Parse(reader)
			if err != nil {
				log.Debug("Parse error:", err)
				return
			}
			requestChan <- &Request{Args: val}
		}
	}()
	return requestChan
}
