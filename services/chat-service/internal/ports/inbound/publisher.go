package inbound_port

import (
	"context"

	"github.com/wnmay/horo/shared/contract"
)

type MessagePublisher interface {
	Publish(ctx context.Context, routingKey string, message contract.AmqpMessage) error
}
