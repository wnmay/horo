package message

import (
	"log"

	"github.com/wnmay/horo/shared/message"
)

var MQ *message.RabbitMQ

// InitRabbitMQ connects to RabbitMQ and returns a ready instance.
func InitRabbitMQ(uri string) *message.RabbitMQ {
	rmq, err := message.NewRabbitMQ(uri)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}

	log.Println("RabbitMQ connected and exchanges/queues declared successfully.")
	MQ = rmq
	return rmq
}
