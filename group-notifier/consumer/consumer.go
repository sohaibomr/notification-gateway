package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sohaibomr/notification-gateway/common/models"
	"github.com/sohaibomr/notification-gateway/common/util"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	groupConsumer sarama.ConsumerGroup
}

func NewConsumerGroup(kafkaBroker []string, groupID string) *Consumer {

	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0
	group, err := sarama.NewConsumerGroup(kafkaBroker, groupID, config)
	if err != nil {
		panic(err)
	}

	return &Consumer{groupConsumer: group}
}

func (c *Consumer) Consume() {

	go func() {
		for err := range c.groupConsumer.Errors() {
			panic(err)
		}
	}()

	func() {
		ctx := context.Background()
		for {
			topics := []string{"group"}
			err := c.groupConsumer.Consume(ctx, topics, c)
			if err != nil {
				fmt.Printf("kafka consume failed: %v, sleeping and retry in a moment\n", err)
				time.Sleep(time.Second)
			}
		}
	}()
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	brokers := []string{"localhost:9092"}
	producer := util.NewKafkaProducer(brokers)
	defer producer.Close()
	for msg := range claim.Messages() {
		jsonMap := make(map[string]interface{})
		json.Unmarshal(msg.Value, &jsonMap)
		groupID := jsonMap["groupId"].(string)
		if util.GroupExist(groupID) {
			for _, user := range models.GroupUsers[groupID] {
				var notification models.UserMsg
				notification.Message = jsonMap["message"].(string)
				notification.SendVia = jsonMap["sendVia"].(string)
				notification.UserDetail = user
				util.NotificationForwarder(&notification, producer) // sends notifiation to sms or push chanel

			}
		}
		defer sess.MarkMessage(msg, "")

	}
	return nil
}
