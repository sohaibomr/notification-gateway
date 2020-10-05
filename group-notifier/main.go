package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/sohaibomr/notification-gateway/common/util"
	"github.com/sohaibomr/notification-gateway/group-notofier/consumer"
	"github.com/spf13/viper"
)

const (
	kafkaGroupIDEnv    = "GROUP_ID"
	groupTopicNamesEnv = "GROUP_TOPIC_NAMES"
	kafkaBrokerAddrEnv = "KAFKA_BROKER_ADDR"
	smsMutexKeyEnv     = "SMS_MUTEX_KEY"
)

var (
	kafkaGroupID    string
	groupTopicNames []string
	kafkaBrokerAddr []string
	rL              *util.RateLimiter
	smsMutexKey     string
)

func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)
	viper.SetDefault(kafkaGroupIDEnv, "groupNotifications")
	viper.SetDefault(groupTopicNamesEnv, []string{"group"})
	viper.SetDefault(kafkaBrokerAddrEnv, []string{"localhost:9092"})
	viper.SetDefault(smsMutexKeyEnv, "sms-mutex")

	// get and set var from env
	viper.BindEnv(kafkaGroupIDEnv)
	viper.BindEnv(groupTopicNamesEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)
	viper.BindEnv(smsMutexKeyEnv)

	kafkaGroupID = viper.GetString(kafkaGroupIDEnv)
	groupTopicNames = viper.GetStringSlice(groupTopicNamesEnv)
	kafkaBrokerAddr = viper.GetStringSlice(kafkaBrokerAddrEnv)
	smsMutexKey = viper.GetString(smsMutexKeyEnv)

	backgroundCtx := context.Background()
	ctx, cancel := context.WithCancel(backgroundCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		defer func() {
			signal.Stop(c)
			cancel()
		}()
		log.Println("starting stop channel routine")
		<-c
		log.Println("stop channel called")
		cancel()
	}()
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "localhost:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)
	mutex := rs.NewMutex(smsMutexKey)

	rL = util.NewRateLimiter(ctx, pool, mutex)
}
func init() {
	setup()
}
func main() {
	consumer := consumer.NewConsumerGroup(kafkaBrokerAddr, kafkaGroupID, groupTopicNames, rL)
	consumer.Consume()
}
