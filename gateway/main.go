package main

// Package classification Swvl notification-gateway API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v2
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
// swagger:meta
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
	viper.SetDefault(serverIPAndPortEnv, ":8082")
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

	// swagger:operation POST /notification/group notification postGroupNotification
	// ---
	// summary: Create Group Notification.
	// description: Group Notification.
	// tags:
	// - Notification Gateway
	// parameters:
	// - name: acount
	//   in: body
	//   description: req body
	//   schema:
	//     "$ref": "#/definitions/groupNotificationRequest"
	// responses:
	//   "200":
	//     description: Group notification object
	//     schema:
	//       "$ref": "#/definitions/groupNotificationRequest"
	r.POST("/notification/group", apiContext.GroupNotification)

	// swagger:operation POST /notification/user notification postUserNotification
	// ---
	// summary: Create User Notification.
	// description: User Notification.
	// tags:
	// - Notification Gateway
	// parameters:
	// - name: acount
	//   in: body
	//   description: req body
	//   schema:
	//     "$ref": "#/definitions/userNotificationRequest"
	// responses:
	//   "200":
	//     description: User notification object
	//     schema:
	//       "$ref": "#/definitions/userNotificationRequest"
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
