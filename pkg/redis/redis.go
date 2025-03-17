package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}

	log.Fatalf("redis")
	return client
}
