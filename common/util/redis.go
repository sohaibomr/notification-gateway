package util

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
)

// RateLimiter redis client for distributes sync locking
type RateLimiter struct {
	Context     context.Context
	RedisPool   redsyncredis.Pool
	RedisMutex  *redsync.Mutex
	StopChannel chan os.Signal
}

func NewRateLimiter(ctx context.Context, redisPool redsyncredis.Pool, redisMutex *redsync.Mutex) *RateLimiter {
	return &RateLimiter{Context: ctx, RedisPool: redisPool, RedisMutex: redisMutex}
}

func GetSmsLimit(conn redsyncredis.Conn) int64 {
	limit, err := conn.Get("smsLimit")
	if err != nil {
		log.Println(err)
	}
	log.Println("SMS Limit", limit)
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return limitInt
}
