package benchmark

import (
	"github.com/go-redis/redis/v8"
	"testing"
)

import (
	"context"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func BenchmarkRedisSet(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		err := redisClient.Set(ctx, "key", "value", 0).Err()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedisGet(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		val, err := redisClient.Get(ctx, "key").Result()
		if err != nil {
			b.Fatal(err)
		}
		if val != "value" {
			b.Fatalf("Expected 'value', got '%s'", val)
		}
	}
}

func BenchmarkRedisIncr(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		err := redisClient.Incr(ctx, "counter").Err()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRedisDel(b *testing.B) {
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		err := redisClient.Del(ctx, "key").Err()
		if err != nil {
			b.Fatal(err)
		}
	}
}
