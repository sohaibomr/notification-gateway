package consumer

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/sohaibomr/notification-gateway/common/models"
	"github.com/sohaibomr/notification-gateway/common/util"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	groupConsumer sarama.ConsumerGroup
	topicNames    []string
	rateLimiter   *util.RateLimiter
}

func NewConsumerGroup(kafkaBroker []string, groupID string, topicNames []string, rl *util.RateLimiter) *Consumer {

	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0
	group, err := sarama.NewConsumerGroup(kafkaBroker, groupID, config)
	if err != nil {
		panic(err)
	}

	return &Consumer{groupConsumer: group, topicNames: topicNames, rateLimiter: rl}
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
			err := c.groupConsumer.Consume(ctx, c.topicNames, c)
			if err != nil {
				log.Printf("kafka consume failed: %v, sleeping and retry in a moment\n", err)
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
	conn, err := c.rateLimiter.RedisPool.Get(c.rateLimiter.Context)
	defer conn.Close()
	if err != nil {
		log.Println(err)
	}
	var smsLimit int64
	for msg := range claim.Messages() {
		jsonMap := make(map[string]interface{})
		json.Unmarshal(msg.Value, &jsonMap)
		sendVia := jsonMap["sendVia"].(string)
		groupID := jsonMap["groupId"].(string)
		if util.GroupExist(groupID) {
			for _, user := range models.GroupUsers[groupID] {
				if sendVia == "sms" {
					// keep retrying after every second to check if sms limit is reset
					for {
						if err := c.rateLimiter.RedisMutex.Lock(); err != nil {
							panic(err)
						}
						smsLimit = util.GetSmsLimit(conn)
						log.Println("Current sms limit remaining for current minute", smsLimit)
						// Release the lock so other processes or threads can obtain a lock.
						if ok, err := c.rateLimiter.RedisMutex.Unlock(); !ok || err != nil {
							log.Fatal("unlock failed")
						}
						if smsLimit > 0 {
							break
						}
						time.Sleep(1 * time.Second)

					}

					if err := c.rateLimiter.RedisMutex.Lock(); err != nil {
						panic(err)
					}

					var notification models.UserMsg
					notification.Message = jsonMap["message"].(string)
					notification.SendVia = sendVia
					notification.UserDetail = user
					util.NotificationForwarder(&notification, producer) // sends notifiation to sms or push chanel
					smsLimit--
					ok, err := conn.Set("smsLimit", strconv.Itoa(int(smsLimit)))
					if !ok || err != nil {
						log.Println(err)
					}
					// Release the lock so other processes or threads can obtain a lock.
					if ok, err := c.rateLimiter.RedisMutex.Unlock(); !ok || err != nil {
						log.Fatal("unlock failed")
					}

				} else if sendVia == "push" {
					var notification models.UserMsg
					notification.Message = jsonMap["message"].(string)
					notification.SendVia = sendVia
					notification.UserDetail = user
					util.NotificationForwarder(&notification, producer) // sends notifiation to sms or push chanel
				}
			}

		}
		sess.MarkMessage(msg, "")
	}

	return nil
}
