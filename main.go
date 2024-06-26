package main

import (
	"github.com/noexcs/redis-go/config"
	"github.com/noexcs/redis-go/redis/handler"
	"github.com/noexcs/redis-go/tcp"
)

func main() {
	// find config file "redis.conf"
	config.Setup()
	tcp.ListenAndServeWithSignal(config.Properties, handler.New())
}
