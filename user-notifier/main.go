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
	"github.com/sohaibomr/notification-gateway/user-notofier/consumer"
	"github.com/spf13/viper"
)

const (
	kafkaUserGroupIDEnv = "USER_GROUP_ID"
	kafkaBrokerAddrEnv  = "KAFKA_BROKER_ADDR"
	userTopicNamesEnv   = "USER_TOPIC_NAMES"
	smsMutexKeyEnv      = "SMS_MUTEX_KEY"
)

var (
	kafkaGroupID    string
	userTopicNames  []string
	kafkaBrokerAddr []string
	rL              *util.RateLimiter
	smsMutexKey     string
)

func init() {
	setup()
}
func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)

	viper.SetDefault(kafkaUserGroupIDEnv, "userNotifications")
	viper.SetDefault(userTopicNamesEnv, []string{"user"})
	viper.SetDefault(kafkaBrokerAddrEnv, []string{"localhost:9092"})
	viper.SetDefault(smsMutexKeyEnv, "sms-mutex")
	// get and set var from env
	viper.BindEnv(kafkaUserGroupIDEnv)
	viper.BindEnv(userTopicNamesEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)
	viper.BindEnv(smsMutexKeyEnv)

	kafkaGroupID = viper.GetString(kafkaUserGroupIDEnv)
	userTopicNames = viper.GetStringSlice(userTopicNamesEnv)
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
func main() {
	consumer := consumer.NewUserConsumerGroup(kafkaBrokerAddr, kafkaGroupID, userTopicNames, rL)
	consumer.Consume()
}
