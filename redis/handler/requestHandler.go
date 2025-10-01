package handler

import (
	"fmt"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/simpleDB"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/client"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"reflect"
	"strings"

	"net"
	"sync"
)

type RequestHandler struct {
	activeConn sync.Map
	Proceeding sync.WaitGroup
	db         database.DB
}

func New() *RequestHandler {
	return &RequestHandler{db: simpleDB.NewSingleDB()}
}

func (h *RequestHandler) Handle(conn net.Conn) {

	// 解析客户端请求
	requestChan := parser.ParseIncomeStream(conn)

	// 创建客户端实例，存入map中
	clientInst := client.New(conn)

	h.activeConn.Store(clientInst, struct{}{})
	// channel被关闭(或连接断开)，删除客户端
	defer h.CloseClient(clientInst)

	// 处理请求
	for request := range requestChan {
		//parseResult(request)
		var response resp2.RespType
		if request.Err != nil {
			break
		} else {
			response = HandleCommand(clientInst, request.Args, h.db)
		}

		if response == nil {
			response = &resp2.SimpleString{Data: "OK"}
		}

		err := clientInst.Write(response.ToBytes())
		if err != nil {
			break
		}
	}
	log.Debug("Client " + conn.RemoteAddr().String() + " disconnected.")
}

func parseResult(request *parser.Request) {

	// 打印request类型
	if request != nil {
		outType := reflect.TypeOf(request.Args)
		log.Debug(fmt.Sprintf("request type: %s", outType.Name()))

		array, ok := request.Args.(*resp2.Array)
		if ok {
			var builder strings.Builder
			for i := 0; i < array.Length; i++ {
				//[*resp2.RespType, *resp2.RespType, *resp2.RespType, *resp2.RespType, ]
				builder.WriteString(fmt.Sprintf("%s, ", reflect.TypeOf((*array).Data[i]).Name()))
			}
			log.Debug(builder.String())
		}

		if request.Args != nil {
			m := fmt.Sprintf("%q", request.Args.String())
			log.Debug("Receive request: ", m)
		}
	}
}

// Close established connections.
func (h *RequestHandler) Close() {
	h.activeConn.Range(func(key, value any) bool {
		c, ok := key.(*client.Client)
		if ok {
			c.Close()
			h.Proceeding.Done()
		}
		return true
	})
}

func (h *RequestHandler) CloseClient(c *client.Client) {
	h.activeConn.Delete(c)
	c.Close()
}
