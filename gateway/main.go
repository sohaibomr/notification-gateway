package main

import (
	"github.com/Shopify/sarama"
	"github.com/sohaibomr/notification-gateway/common/util"
	api "github.com/sohaibomr/notification-gateway/gateway/api/v1"

	"log"

	"github.com/gin-gonic/gin"
)

var (
	kafkaProducer sarama.AsyncProducer
	apiContext    *api.APIContext
)

func setup() {
	// getting of env vars here
	kafkaProducer = util.NewKafkaProducer([]string{"localhost:9092"})
}

func startAPI() {
	apiContext = api.NewAPIContext(kafkaProducer, "group", "user")
	r := gin.Default()
	r.POST("/notification/group", apiContext.GroupNotification)

	r.POST("/notification/personalized", apiContext.PersonalizedNotification)
	// go func() {
	// 	for err := range kafkaProducer.Errors() {
	// 		log.Println("Failed to write access log entry:", err)
	// 	}
	// }()
	log.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080

}

func main() {
	setup()
	startAPI()
}
