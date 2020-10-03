package api

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

type personalizedNotificationRequest struct {
	ID         string   `json:"userId"`
	Type       string   `json:"type"` //group or custom/personalized
	SendVia    string   `json:"sendVia"`
	Message    string   `json:"message"`
	Category   string   `json:"category"` // promo code, destination reached, captain waiting etc
	Tags       []string `json:"tags"`
	TemplateID string   `json:"templateID"`
}

// PersonalizedNotification pushes noti to kafka topic
func (ac *APIContext) PersonalizedNotification(c *gin.Context) {

	var req personalizedNotificationRequest
	c.BindJSON(&req)
	fmt.Println(req)
	payload, err := json.Marshal(req)
	if err != nil {
		fmt.Println(err)
	}

	// TODO: return error, add schema validation
	ac.kafkaProducer.Input() <- &sarama.ProducerMessage{Topic: ac.userTopic, Value: sarama.ByteEncoder(payload)}
	c.JSON(200, req)

}
