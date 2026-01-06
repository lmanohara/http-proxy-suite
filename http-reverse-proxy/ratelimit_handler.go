package main

import (
	"context"
	"fmt"
	"log"
)

func setRateLimit(ctx context.Context) {
	client := NewRedisClient("redis-server:6379", "", 0)
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	fmt.Println("Connected to the Redis server")
}

func isRateLimited() bool {
	return false
}
