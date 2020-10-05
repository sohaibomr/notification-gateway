module github.com/sohaibomr/notification-gateway/gateway

go 1.13

require (
	github.com/Shopify/sarama v1.27.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/sohaibomr/notification-gateway/common v0.0.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.1.9 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
	golang.org/x/sys v0.0.0-20200929083018-4d22bbb62b3c // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/sohaibomr/notification-gateway/common => ../common
