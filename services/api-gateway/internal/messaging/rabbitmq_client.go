package messaging

import (
	"fmt"
	"log"

	"github.com/wnmay/horo/shared/message"
)

type RabbitMQClient struct {
	rmq *message.RabbitMQ
}

func NewRabbitMQClient(uri string) (*RabbitMQClient, error) {
	rmq, err := message.NewRabbitMQ(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ connection: %v", err)
	}

	log.Println("RabbitMQ client initialized successfully")

	return &RabbitMQClient{
		rmq: rmq,
	}, nil
}

func (r *RabbitMQClient) GetRabbitMQ() *message.RabbitMQ {
	return r.rmq
}

func (r *RabbitMQClient) Close() {
	if r.rmq != nil {
		r.rmq.Close()
		log.Println("RabbitMQ client closed")
	}
}
