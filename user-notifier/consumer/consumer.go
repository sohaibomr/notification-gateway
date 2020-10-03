package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sohaibomr/notification-gateway/common/models"
	"github.com/sohaibomr/notification-gateway/common/util"
)

// UserConsumer ..
type UserConsumer struct {
	group      sarama.ConsumerGroup
	topicNames []string
}

func NewUserConsumerGroup(kafkaBroker []string, groupID string, topicNames []string) *UserConsumer {

	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0
	group, err := sarama.NewConsumerGroup(kafkaBroker, groupID, config)
	if err != nil {
		panic(err)
	}

	return &UserConsumer{group: group, topicNames: topicNames}
}

func (c *UserConsumer) Consume() {

	go func() {
		for err := range c.group.Errors() {
			panic(err)
		}
	}()

	func() {
		ctx := context.Background()
		for {
			err := c.group.Consume(ctx, c.topicNames, c)
			if err != nil {
				fmt.Printf("kafka consume failed: %v, sleeping and retry in a moment\n", err)
				time.Sleep(time.Second)
			}
		}
	}()
}

func (c *UserConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *UserConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *UserConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	brokers := []string{"localhost:9092"}
	producer := util.NewKafkaProducer(brokers)
	defer producer.Close()
	for msg := range claim.Messages() {
		fmt.Printf("consumed a message: %v\n", string(msg.Value))
		jsonMap := make(map[string]interface{})
		json.Unmarshal(msg.Value, &jsonMap)
		var notification models.UserMsg
		notification.Message = msg["message"].(string)
		notification.SendVia = msg["sendVia"].(string)
		userID := msg["userId"].(string)
		if util.UserExist(userID) {
			notification.UserDetail = models.UsersMap[userID]
		}
		util.NoificationForwarder(&notification, producer) // sends notifiation to sms or push chanel
		sess.MarkMessage(msg, "")
	}
	return nil
}
