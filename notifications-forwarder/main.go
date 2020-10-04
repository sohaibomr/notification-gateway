package main

import (
	"log"
	"sync"
	"time"

	"github.com/sohaibomr/notification-gateway/notifications-forwarder/consumer"
	"github.com/spf13/viper"
)

const (
	smsGroupIDEnv      = "SMS_GROUP_ID"
	pushGroupIDEnv     = "PUSH_GROUP_ID"
	smsTopicNamesEnv   = "SMS_TOPICS"
	pushTopicNamesEnv  = "PUSH_TOPICS"
	kafkaBrokerAddrEnv = "KAFKA_BROKER_ADDR"
)

var (
	smsGroupID              string
	pushNotificationGroupID string
	smsTopicNames           []string
	pushTopicNames          []string
	kafkabrokerAddr         []string
)

func init() {
	setup()
}
func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)
	viper.SetDefault(smsGroupIDEnv, "smsNotifications")
	viper.SetDefault(pushGroupIDEnv, "pushNotifications")
	viper.SetDefault(smsTopicNamesEnv, []string{"sms"})
	viper.SetDefault(pushTopicNamesEnv, []string{"push"})
	viper.SetDefault(kafkaBrokerAddrEnv, []string{"localhost:9092"})

	// get and set var from env
	viper.BindEnv(smsGroupIDEnv)
	viper.BindEnv(pushGroupIDEnv)
	viper.BindEnv(smsTopicNamesEnv)
	viper.BindEnv(pushTopicNamesEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)

	smsGroupID = viper.GetString(smsGroupIDEnv)
	pushNotificationGroupID = viper.GetString(pushGroupIDEnv)
	smsTopicNames = viper.GetStringSlice(smsTopicNamesEnv)
	pushTopicNames = viper.GetStringSlice(pushTopicNamesEnv)
	kafkabrokerAddr = viper.GetStringSlice(kafkaBrokerAddrEnv)
}
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	pushConsumer := consumer.NewPushNotificationConsumerGroup(kafkabrokerAddr, pushNotificationGroupID, pushTopicNames)
	go func() {
		pushConsumer.ConsumePushNotification()
		defer wg.Done()
	}()
	smsConsumer := consumer.NewSmsConsumerGroup(kafkabrokerAddr, smsGroupID, smsTopicNames)
	go func() {
		smsConsumer.ConsumeSms()
		defer wg.Done()
	}()
	wg.Wait()
}
