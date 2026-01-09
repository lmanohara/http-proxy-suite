package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func setRateLimit(ctx context.Context) *redis.Client {
	client := NewRedisClient("redis-server:6379", "", 0)
	// defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Could not connect to Redis:", err)
	}

	// set the radis cache to keep permits per minutes

	fmt.Println("Connected to the Redis server")

	return client
}

func isRateLimited(clientId string, ctx context.Context, client *redis.Client) bool {

	keyCount := fmt.Sprintf("rate_limit:%s:count", clientId)
	keyLastRefill := fmt.Sprintf("rate_limit:%s:last_refill", clientId)

	currentTime := time.Now().UnixMilli()

	pipline := client.TxPipeline()

	getKeyCount := pipline.Get(ctx, keyCount)
	getKeyLastRefill := pipline.Get(ctx, keyLastRefill)

	_, err := pipline.Exec(ctx)
	if err != nil {
		log.Println("Error executing pipeline:", err)
	}

	strRequestCount, err := getKeyCount.Result()
	requestCount := 0

	if err == nil {
		if rc, err := strconv.Atoi(strRequestCount); err == nil {
			requestCount = rc
		}
	}
	strLastLeakTime, err := getKeyLastRefill.Result()

	lastLeakTime := currentTime

	if err == nil {
		if lt, err := strconv.ParseInt(strLastLeakTime, 10, 64); err == nil {
			lastLeakTime = lt
		}
	}

	elapsedTimeInMs := currentTime - lastLeakTime
	elapsedTimeInSeconds := float64(elapsedTimeInMs) / 1000.0

	requestToLeak := int(elapsedTimeInSeconds * 1.0) // 1 request per second

	requestCount -= requestToLeak
	if requestToLeak < 0 {
		requestCount = 0
	}
	isAllowed := requestCount < 5 // max 5 requests in the bucket

	if isAllowed {
		requestCount++
	}

	pipline = client.TxPipeline()

	pipline.Set(ctx, keyCount, requestCount, 0)
	pipline.Set(ctx, keyLastRefill, currentTime, 0)

	pipline.Exec(ctx)

	fmt.Printf("Client %s - Request Count: %d and Last refill time: %d, Allowed: %v\n", clientId, requestCount, lastLeakTime, isAllowed)
	return isAllowed
}
