package main

import (
	"log"
	"time"

	"github.com/sohaibomr/notification-gateway/group-notofier/consumer"
	"github.com/spf13/viper"
)

const (
	kafkaGroupIDEnv    = "GROUP_ID"
	groupTopicNamesEnv = "GROUP_TOPIC_NAMES"
	kafkaBrokerAddrEnv = "KAFKA_BROKER_ADDR"
)

var (
	kafkaGroupID    string
	groupTopicNames []string
	kafkaBrokerAddr []string
)

func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)
	viper.SetDefault(kafkaGroupIDEnv, "groupNotifications")
	viper.SetDefault(groupTopicNamesEnv, []string{"group"})
	viper.SetDefault(kafkaBrokerAddrEnv, []string{"localhost:9092"})

	// get and set var from env
	viper.BindEnv(kafkaGroupIDEnv)
	viper.BindEnv(groupTopicNamesEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)

	kafkaGroupID = viper.GetString(kafkaGroupIDEnv)
	groupTopicNames = viper.GetStringSlice(groupTopicNamesEnv)
	kafkaBrokerAddr = viper.GetStringSlice(kafkaBrokerAddrEnv)
}
func init() {
	setup()
}
func main() {
	consumer := consumer.NewConsumerGroup(kafkaBrokerAddr, kafkaGroupID, groupTopicNames)
	consumer.Consume()
}
