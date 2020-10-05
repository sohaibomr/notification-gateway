package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/sohaibomr/notification-gateway/common/models"
	"github.com/sohaibomr/notification-gateway/common/util"

	"github.com/gin-gonic/gin"
)

// swagger:model userNotificationRequest
type personalizedNotificationRequest struct {
	UserID string `json:"userId" binding:"required"`
	models.NotificationRequest
}

func Simple(verr validator.ValidationErrors) map[string]string {
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

// PersonalizedNotification pushes noti to kafka topic
func (ac *APIContext) PersonalizedNotification(c *gin.Context) {

	var req personalizedNotificationRequest
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

	ac.kafkaProducer.Input() <- &sarama.ProducerMessage{Topic: ac.userTopic, Value: sarama.ByteEncoder(payload)}
	c.JSON(http.StatusOK, req)

}
