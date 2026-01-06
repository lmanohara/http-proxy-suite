package main

import (
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient(addr string, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}
