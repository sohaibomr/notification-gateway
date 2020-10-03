package api

import "github.com/Shopify/sarama"

type APIContext struct {
	kafkaProducer sarama.AsyncProducer
	groupTopic    string
	userTopic     string
}

// NewAPiContext returns api context object to use within objects
func NewAPIContext(kafkaProducer sarama.AsyncProducer, groupTopic string, userTopic string) *APIContext {
	return &APIContext{kafkaProducer: kafkaProducer, groupTopic: groupTopic, userTopic: userTopic}

}
