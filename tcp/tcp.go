package tcp

import (
	"context"
	"fmt"
	"github.com/noexcs/redis-go/command"
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/database"
	"github.com/noexcs/redis-go/database/simpleDB"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/client"
	"github.com/noexcs/redis-go/redis/handler"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"net"
	"strings"
	"sync"
	"time"
)

type dbRequest struct {
	Client  *client.Client
	Args    resp.RespValue
	Result  chan resp.RespValue
	IsWrite bool // 标记是否为写命令
}

type Server struct {
	db         database.DB
	aofHandler *database.AofHandler
	listener   net.Listener
	running    bool
	mutex      sync.Mutex
	wg         sync.WaitGroup

	activeConn sync.Map
	Proceeding sync.WaitGroup

	// 数据库请求通道
	dbChan chan *dbRequest
}

func NewServer() *Server {
	db := simpleDB.NewGoMapDB()

	var aofHandler *database.AofHandler
	if config.Properties.AppendOnly {
		var err error
		aofHandler, err = database.NewAofHandler(db, config.Properties.AppendFilename)
		if err != nil {
			log.Error("Failed to initialize AOF handler: %v", err)
		}
	}

	return &Server{
		db:         db,
		aofHandler: aofHandler,
		dbChan:     make(chan *dbRequest, 1000), // 缓冲1000个请求
	}
}

func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	s.listener = listener
	s.running = true
	s.mutex.Unlock()

	log.Info(fmt.Sprintf("Server started on %s", address))

	// 启动数据库处理 goroutine
	go s.processDBRequests()

	// 启动定期清理过期键的任务
	go s.startPeriodicCleanup()

	for {
		s.mutex.Lock()
		running := s.running
		s.mutex.Unlock()

		if !running {
			break
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Debug("Accept error: " + err.Error())
			continue
		}
		log.Debug("New connection from " + conn.RemoteAddr().String())
		s.wg.Add(1)
		go s.Handle(conn)
	}

	return nil
}

func (s *Server) Handle(conn net.Conn) {

	// 解析客户端请求
	requestChan := parser.ParseIncomeStream(conn)

	// 创建客户端实例，存入map中
	clientInst := client.New(conn)

	s.activeConn.Store(clientInst, struct{}{})
	// channel被关闭(或连接断开)，删除客户端
	defer s.CloseClient(clientInst)

	// 处理请求
	for request := range requestChan {
		//parseResult(request)
		var response resp.RespValue
		if request.Err != nil {
			break
		}
		log.Debug("Client command: ", request.Args.String())
		// 验证命令
		response = handler.ValidateCommand(clientInst, request.Args)
		if response == nil {
			// 检查是否为写命令
			isWrite := false
			if array, ok := request.Args.(*resp2.Array); ok && len(array.Data) > 0 {
				commandName := array.Data[0].String()
				if cmd, exists := command.CmdTable[strings.ToUpper(commandName)]; exists {
					isWrite = (cmd.Flags & command.FlagWrite) != 0
				}
			}

			// 将请求发送到数据库处理 goroutine
			resultChan := make(chan resp.RespValue, 1)
			s.dbChan <- &dbRequest{
				Client:  clientInst,
				Args:    request.Args,
				Result:  resultChan,
				IsWrite: isWrite,
			}

			// 等待结果
			response = <-resultChan
			close(resultChan)
		}
		if response == nil {
			response = &resp2.SimpleString{Data: "OK"}
		}

		err := clientInst.Write(response.ToBytes())
		if err != nil {
			log.Debug("Write response error: " + err.Error())
			break
		}
	}
	log.Debug("Client " + conn.RemoteAddr().String() + " disconnected.")
}

// processDBRequests 在单独的 goroutine 中处理所有数据库请求
func (s *Server) processDBRequests() {
	for req := range s.dbChan {
		response := handler.ExecCommand(req.Client, req.Args, s.db)

		// 如果启用了AOF并且命令是写操作，则将命令添加到AOF
		if s.aofHandler != nil && req.IsWrite {
			// 检查命令是否为数组类型
			if array, ok := req.Args.(*resp2.Array); ok {
				s.aofHandler.AddAof(array.Data)
			}
		}

		req.Result <- response
		s.db.RandomExpiredKeys()
	}
}

func (s *Server) CloseClient(c *client.Client) {
	s.activeConn.Delete(c)
	c.Close()
}

// startPeriodicCleanup 启动定期清理任务
func (s *Server) startPeriodicCleanup() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		s.mutex.Lock()
		running := s.running
		s.mutex.Unlock()

		if !running {
			break
		}

		select {
		case <-ticker.C:
			// 定期执行全面的过期键清理
			s.db.DeleteExpiredKeys()
		case <-time.After(10 * time.Second):
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.mutex.Lock()
	if !s.running {
		s.mutex.Unlock()
		return nil
	}

	s.running = false
	if s.listener != nil {
		_ = s.listener.Close()
	}
	s.mutex.Unlock()

	// 等待所有连接处理完成
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("Server shutdown complete")
		// 关闭数据库请求通道
		close(s.dbChan)

		// 关闭AOF处理器
		if s.aofHandler != nil {
			s.aofHandler.Close()
		}
	case <-ctx.Done():
		log.Info("Server shutdown timeout")
		return ctx.Err()
	}

	return nil
}
