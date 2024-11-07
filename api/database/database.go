package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateClient(dbNo int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("DB_ADDRESS"),
		Password:     os.Getenv("DB_PASSWORD"),
		DB:           dbNo,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	// Test connection with retry
	for i := 0; i < 3; i++ {
		_, err := rdb.Ping(Ctx).Result()
		if err == nil {
			return rdb, nil
		}
		fmt.Printf("Failed to connect to Redis (attempt %d/3): %v\n", i+1, err)
		time.Sleep(time.Second * 2)
	}

	return nil, fmt.Errorf("failed to connect to Redis after 3 attempts")
}
