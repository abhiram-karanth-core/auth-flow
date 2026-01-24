package redisclient

import (
	"context"
	// "crypto/tls"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func New() *redis.Client {
	redisURL := os.Getenv("REDIS_ADDR")
	if redisURL == "" {
		log.Fatal("REDIS_ADDR not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Invalid REDIS_ADDR: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}

	return client
}