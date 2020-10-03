package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// SmsConsumer ..
type SmsConsumer struct {
	group      sarama.ConsumerGroup
	topicNames []string
}

func NewSmsConsumerGroup(kafkaBroker []string, groupID string, topicNames []string) *SmsConsumer {

	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0
	group, err := sarama.NewConsumerGroup(kafkaBroker, groupID, config)
	if err != nil {
		panic(err)
	}

	return &SmsConsumer{group: group, topicNames: topicNames}
}

func (c *SmsConsumer) ConsumeSms() {

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

func (c *SmsConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *SmsConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *SmsConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("consumed a message from SMS channel: %v\n", string(msg.Value))
		sess.MarkMessage(msg, "")
	}
	return nil
}
