package main

import (
	"log"
	"time"

	"github.com/sohaibomr/notification-gateway/user-notofier/consumer"
	"github.com/spf13/viper"
)

const (
	kafkaUserGroupIDEnv = "USER_GROUP_ID"
	kafkaBrokerAddrEnv  = "KAFKA_BROKER_ADDR"
	userTopicNamesEnv   = "USER_TOPIC_NAMES"
)

var (
	kafkaGroupID    string
	userTopicNames  []string
	kafkaBrokerAddr []string
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

	// get and set var from env
	viper.BindEnv(kafkaUserGroupIDEnv)
	viper.BindEnv(userTopicNamesEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)

	kafkaGroupID = viper.GetString(kafkaUserGroupIDEnv)
	userTopicNames = viper.GetStringSlice(userTopicNamesEnv)
	kafkaBrokerAddr = viper.GetStringSlice(kafkaBrokerAddrEnv)

	//get consumer config
	kafkaGroupID = "userNotifications"
}
func main() {
	kafkaBrokers := []string{"localhost:9092"}
	consumer := consumer.NewUserConsumerGroup(kafkaBrokers, kafkaGroupID, userTopicNames)
	consumer.Consume()
}
