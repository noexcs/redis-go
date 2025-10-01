package tcp

import (
	"context"
	"fmt"
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/database/simpleDB"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/client"
	"github.com/noexcs/redis-go/redis/handler"
	"github.com/noexcs/redis-go/redis/parser"
	"github.com/noexcs/redis-go/redis/parser/resp2"
	"net"
	"sync"
	"time"
)

type Server struct {
	db       *simpleDB.SingleDB
	listener net.Listener
	running  bool
	mutex    sync.Mutex
	wg       sync.WaitGroup

	activeConn sync.Map
	Proceeding sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		db: simpleDB.NewSingleDB(),
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
		var response resp2.RespType
		if request.Err != nil {
			break
		} else {
			response = handler.HandleCommand(clientInst, request.Args, s.db)
		}

		if response == nil {
			response = &resp2.SimpleString{Data: "OK"}
		}

		err := clientInst.Write(response.ToBytes())
		if err != nil {
			log.Debug("Write response error: " + err.Error())
			break
		}
		s.db.RandomExpiredKeys()
	}
	log.Debug("Client " + conn.RemoteAddr().String() + " disconnected.")
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
	case <-ctx.Done():
		log.Info("Server shutdown timeout")
		return ctx.Err()
	}

	return nil
}
