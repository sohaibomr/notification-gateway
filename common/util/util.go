package util

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sohaibomr/notification-gateway/common/models"

	"github.com/Shopify/sarama"
)

func NewKafkaProducer(brokerList []string) sarama.AsyncProducer {
	// brokerList := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func UserExist(id string) bool {
	_, ok := models.UsersMap[id]
	return ok
}
func GroupExist(id string) bool {
	_, ok := models.GroupUsers[id]
	return ok
}
func sendSMS(notification *models.UserMsg, producer sarama.AsyncProducer) {

	payload, err := json.Marshal(notification)
	if err != nil {
		log.Println(err)
	}
	producer.Input() <- &sarama.ProducerMessage{Topic: "sms", Value: sarama.ByteEncoder(payload)}

}
func sendPushNotification(notification *models.UserMsg, producer sarama.AsyncProducer) {

	payload, err := json.Marshal(notification)
	if err != nil {
		log.Println(err)
	}
	producer.Input() <- &sarama.ProducerMessage{Topic: "push", Value: sarama.ByteEncoder(payload)}
}

// MessageForwarder gets the user meta from DB and sends notification to sms or push service
func NotificationForwarder(notification *models.UserMsg, producer sarama.AsyncProducer) {
	log.Println(notification)
	sendVia := notification.SendVia

	if sendVia == "sms" {
		sendSMS(notification, producer)
	} else if sendVia == "push" {
		sendPushNotification(notification, producer)
	}

}

func GinValidationErr(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
