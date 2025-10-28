package inbound_port

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageConsumer interface {
	Consume(queueName string, handler func(context.Context, amqp.Delivery) error) error
	StartListening() error
}
