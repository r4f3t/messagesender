package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
)

// InitializeRedis sets up the connection to the Redis database
func InitializeRedis() *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379",
		Password: "", // No password set
		DB:       0,  // Default DB
	})

	if _, err := redisDB.Ping(ctx).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return nil
	}

	return redisDB
}

// CacheMessage caches the message ID and sending time in Redis
func CacheMessage(messageID uint, redisDB *redis.Client) error {
	cacheKey := fmt.Sprintf("message:%d", messageID)
	cacheValue := time.Now().Format(time.RFC3339)
	return redisDB.Set(ctx, cacheKey, cacheValue, 0).Err()
}
