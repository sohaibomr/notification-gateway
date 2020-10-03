package main

import "github.com/sohaibomr/notification-gateway/group-notofier/consumer"

// add health API
var (
	kafkaGroupID string
)

func setup() {
	//get consumer config
	kafkaGroupID = "groupNotifications"

}
func init() {
	setup()
}
func main() {
	kafkaBrokers := []string{"localhost:9092"}
	consumer := consumer.NewConsumerGroup(kafkaBrokers, kafkaGroupID)
	consumer.Consume()
}
