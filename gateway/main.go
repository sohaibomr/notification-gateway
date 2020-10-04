package main

import (
	"reflect"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/sohaibomr/notification-gateway/common/util"
	api "github.com/sohaibomr/notification-gateway/gateway/api/v1"
	"github.com/spf13/viper"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	groupTopicNameEnv  = "GROUP_TOPIC"
	userTopicNameEnv   = "USER_TOPIC"
	serverIPAndPortEnv = "SERVER_IP_PORT"
	kafkaBrokerAddrEnv = "KAFKA_BROKER_ADDR"
)

var (
	kafkaProducer   sarama.AsyncProducer
	apiContext      *api.APIContext
	groupTopicName  string
	userTopicName   string
	serverIPAndPort string
	kafkaBrokerAddr []string
)

func init() {
	setup()
}
func setup() {
	log.Println("Waiting for kafka to start...")
	time.Sleep(20 * time.Second)
	// set defaults
	viper.SetDefault(groupTopicNameEnv, "group")
	viper.SetDefault(userTopicNameEnv, "user")
	viper.SetDefault(serverIPAndPortEnv, "0.0.0.0:8080")
	viper.SetDefault(kafkaBrokerAddrEnv, []string{"localhost:9092"})

	// get and set var from env
	viper.BindEnv(groupTopicNameEnv)
	viper.BindEnv(userTopicNameEnv)
	viper.BindEnv(serverIPAndPortEnv)
	viper.BindEnv(kafkaBrokerAddrEnv)
	groupTopicName = viper.GetString(groupTopicNameEnv)
	userTopicName = viper.GetString(userTopicNameEnv)
	serverIPAndPort = viper.GetString(serverIPAndPortEnv)
	kafkaBrokerAddr = viper.GetStringSlice(kafkaBrokerAddrEnv)

	kafkaProducer = util.NewKafkaProducer(kafkaBrokerAddr)
	apiContext = api.NewAPIContext(kafkaProducer, groupTopicName, userTopicName)
}

func startAPI() {
	r := gin.Default()
	// register custom tag extractor function for Gin default schema validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	r.POST("/notification/group", apiContext.GroupNotification)

	r.POST("/notification/user", apiContext.PersonalizedNotification)
	go func() {
		for err := range kafkaProducer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()
	log.Fatal(r.Run(serverIPAndPort))

}

func main() {
	startAPI()
}
