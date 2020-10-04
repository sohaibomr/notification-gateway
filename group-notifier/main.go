package main

import (
	"log"
	"time"

	"github.com/sohaibomr/notification-gateway/group-notofier/consumer"
)

// add health API
var (
	kafkaGroupID string
)

func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)
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
