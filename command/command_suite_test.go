package command_test

import (
	"context"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Suite")
}

var ctx = context.Background()
var client *redis.Client

var _ = BeforeSuite(func() {
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6390",
		//Addr:     "localhost:6398",
		//Password: "123456",
		DB: 0,
	})
	client.FlushDB(ctx)
})

var _ = AfterSuite(func() {
	client.FlushDB(ctx)
	// 关闭连接
	client.Close()
})
