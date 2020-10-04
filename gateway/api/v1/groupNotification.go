package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sohaibomr/notification-gateway/common/models"

	"github.com/sohaibomr/notification-gateway/common/util"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

type groupNotificationRequest struct {
	GroupID string `json:"groupId" binding:"required"`
	models.NotificationRequest
}

// Notification pushes noti to kafka topic
func (ac *APIContext) GroupNotification(c *gin.Context) {

	var req groupNotificationRequest
	if err := c.ShouldBind(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": util.GinValidationErr(verr)})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	ac.kafkaProducer.Input() <- &sarama.ProducerMessage{Topic: ac.groupTopic, Value: sarama.ByteEncoder(payload)}
	c.JSON(http.StatusOK, req)

}
