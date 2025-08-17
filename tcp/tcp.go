package tcp

import (
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/redis/handler"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var closeChan = make(chan struct{})

func ListenAndServeWithSignal(config *config.ServerProperties, handler *handler.RequestHandler) {

	// 接收系统信号用于停止程序
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", config.Bind+":"+strconv.FormatInt(int64(config.Port), 10))
	if err != nil {
		log.ForceWithLocation(err)
		return
	}

	log.ForceWithLocation("Listening on " + config.Bind + ":" + strconv.FormatInt(int64(config.Port), 10) + "")

	ListenAndServe(listener, handler, closeChan)
}

func ListenAndServe(listener net.Listener, handler *handler.RequestHandler, closeChan chan struct{}) {

	go func() {
		// Receives System Signal from sigCh
		<-closeChan
		log.ForceWithLocation("\033[04mShutting down...")

		// Stop to accept new connection.
		listener.Close()
		// Close established connections.
		handler.Close()
	}()

	log.ForceWithLocation("The server has been started.")

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		log.WithLocation("Accepted connection from " + conn.RemoteAddr().String() + ".")
		go func(conn2 net.Conn) {
			handler.Proceeding.Add(1)
			defer handler.Proceeding.Done()

			handler.Handle(conn2)
		}(conn)
	}

	handler.Proceeding.Wait()
}

func StopService() {
	closeChan <- struct{}{}
}
