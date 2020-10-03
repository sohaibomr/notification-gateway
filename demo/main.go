package main

import (
	"sync"

	"github.com/sohaibomr/notification-gateway/demo/consumer"
)

var (
	smsGroupID              string
	pushNotificationGroupID string
	smsTopicNames           []string
	pushTopicNames          []string
)

func init() {
	setup()
}
func setup() {
	//get consumer config
	smsGroupID = "smsNotifications"
	pushNotificationGroupID = "pushNotifications"
	smsTopicNames = []string{"sms"}
	pushTopicNames = []string{"push"}
}
func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	kafkaBrokers := []string{"localhost:9092"}
	pushConsumer := consumer.NewPushNotificationConsumerGroup(kafkaBrokers, pushNotificationGroupID, pushTopicNames)
	go func() {
		pushConsumer.ConsumePushNotification()
		defer wg.Done()
	}()
	smsConsumer := consumer.NewSmsConsumerGroup(kafkaBrokers, smsGroupID, smsTopicNames)
	go func() {
		smsConsumer.ConsumeSms()
		defer wg.Done()
	}()
	wg.Wait()
}
