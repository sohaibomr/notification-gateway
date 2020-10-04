module github.com/sohaibomr/notification-gateway/user-notofier

go 1.13

require github.com/Shopify/sarama v1.27.0

require (
	github.com/sohaibomr/notification-gateway/common v0.0.0
	github.com/spf13/viper v1.7.1
)

replace github.com/sohaibomr/notification-gateway/common => ../common
