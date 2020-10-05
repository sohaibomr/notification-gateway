module github.com/sohaibomr/notification-gateway/group-notofier

go 1.13

require github.com/Shopify/sarama v1.27.0

require (
	github.com/go-redis/redis/v8 v8.1.1
	github.com/go-redsync/redsync/v4 v4.0.3
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/sohaibomr/notification-gateway/common v0.0.0
	github.com/spf13/viper v1.7.1
)

replace github.com/sohaibomr/notification-gateway/common => ../common
