package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/spf13/viper"
)

const (
	smsLimitEnv    = "SMS_LIMIT"
	smsMutexEnv    = "SMS_MUTEX"
	smsRedisKeyEnv = "SMS_REDIS_KEY"
	redisAddrEnv   = "REDIS_ADDR"
)

var (
	smsLimit    string
	smsMutex    string
	smsRedisKey string
	redisAddr   string
)

func init() {
	setup()
}

func setup() {
	time.Sleep(20 * time.Second)
	viper.SetDefault(smsLimitEnv, "100")
	viper.SetDefault(smsMutexEnv, "sms-mutex")
	viper.SetDefault(smsRedisKeyEnv, "smsLimit")
	viper.SetDefault(redisAddrEnv, "localhost:6379")

	viper.BindEnv(smsLimitEnv)
	viper.BindEnv(smsMutexEnv)
	viper.BindEnv(smsRedisKeyEnv)
	viper.BindEnv(redisAddrEnv)

	smsLimit = viper.GetString(smsLimitEnv)
	smsMutex = viper.GetString(smsMutexEnv)
	smsRedisKey = viper.GetString(smsRedisKeyEnv)
	redisAddr = viper.GetString(redisAddrEnv)
}

type rateLimiter struct {
	context     context.Context
	redisPool   redsyncredis.Pool
	redisMutex  *redsync.Mutex
	stopChannel chan os.Signal
}

func newRateLimiter(ctx context.Context, redisPool redsyncredis.Pool, redisMutex *redsync.Mutex) *rateLimiter {
	return &rateLimiter{context: ctx, redisPool: redisPool, redisMutex: redisMutex}
}

func (rl *rateLimiter) resetSmsCounter() {
	if err := rl.redisMutex.Lock(); err != nil {
		panic(err)
	}
	conn, err := rl.redisPool.Get(rl.context)
	defer conn.Close()
	if err != nil {
		log.Println(err)
	}
	ok, err := conn.Set(smsRedisKey, smsLimit)
	if !ok || err != nil {
		log.Println(err)
	}
	limit, err := conn.Get("smsLimit")
	if err != nil {
		log.Println(err)
	}
	log.Println("SMS Limit", limit)
	// Release the lock so other processes or threads can obtain a lock.
	if ok, err := rl.redisMutex.Unlock(); !ok || err != nil {
		log.Fatal("unlock failed")
	}

}
func (rl *rateLimiter) resetRateCounter(seconds time.Duration) error {
	ticker := time.Tick(seconds)
	for {
		select {
		case <-rl.context.Done():
			return rl.context.Err()
		case <-ticker:
			log.Println("Resetting sms rate limit counter")
			rl.resetSmsCounter()
			log.Println("SMS rate limit reset complete")

		}
	}

}

func main() {
	backgroundCtx := context.Background()
	ctx, cancel := context.WithCancel(backgroundCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		log.Println("starting stop channel routine")
		<-c
		log.Println("stop channel called")
		cancel()
	}()

	client := goredislib.NewClient(&goredislib.Options{
		Addr: redisAddr,
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)

	mutex := rs.NewMutex(smsMutex)

	rL := newRateLimiter(ctx, pool, mutex)
	rL.resetSmsCounter()
	err := rL.resetRateCounter(60 * time.Second)

	log.Println(err)
}
