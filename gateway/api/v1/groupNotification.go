package api

import (
	"encoding/json"
	"log"

	"github.com/sohaibomr/notification-gateway/common/util"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

type groupNotificationRequest struct {
	Type     string   `json:"type"` //group or custom/personalized
	SendVia  string   `json:"sendVia"`
	Message  string   `json:"message"`
	GroupID  string   `json:"groupId"`
	Category string   `json:"category"` // promo code, destination reached, captain waiting etc
	Tags     []string `json:"tags"`
}

// Notification pushes noti to kafka topic
func (ac *APIContext) GroupNotification(c *gin.Context) {

	var req groupNotificationRequest
	c.BindJSON(&req)
	payload, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	// TODO: return error, add schema validation
	ac.kafkaProducer.Input() <- &sarama.ProducerMessage{Topic: ac.groupTopic, Value: sarama.ByteEncoder(payload)}
	// TODO:gracefully cloe application and close producer
	c.JSON(util.StatusCodeOK, req)

}
