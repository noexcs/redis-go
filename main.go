package main

import (
	"context"
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/log"
	"github.com/noexcs/redis-go/tcp"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// find a config file "redis.conf"
	// 读取配置文件redis.conf
	config.Setup()

	server := tcp.NewServer()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Start the server
	err := server.Start()
	if err != nil {
		log.Info("Failed to start server: " + err.Error())
	}
}
