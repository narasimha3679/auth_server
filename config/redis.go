package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Ping Redis to verify connection
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}
