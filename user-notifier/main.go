package main

import "github.com/sohaibomr/notification-gateway/user-notofier/consumer"

// add health API
var (
	kafkaGroupID string
	topicNames   []string
)

func init() {
	setup()
}
func setup() {
	//get consumer config
	kafkaGroupID = "userNotifications"
	topicNames = []string{"user"}
}
func main() {
	kafkaBrokers := []string{"localhost:9092"}
	consumer := consumer.NewUserConsumerGroup(kafkaBrokers, kafkaGroupID, topicNames)
	consumer.Consume()
}
