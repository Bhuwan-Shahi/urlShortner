package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

// Ctx is the background context used for Redis operations
var Ctx = context.Background()

// CreateClient creates and returns a new Redis client
// dbNo parameter specifies which Redis database to use
func CreateClient(dbNo int) *redis.Client {
	// Create new Redis client with configuration from environment variables
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),  // Fixed environment variable name
		Password: os.Getenv("DB_PASSWORD"), // Fixed environment variable name
		DB:       dbNo,
	})

	// Verify connection
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	return rdb
}
