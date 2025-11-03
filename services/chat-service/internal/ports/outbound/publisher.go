package outbound_port

import (
	"context"

	"github.com/wnmay/horo/shared/contract"
)

type MessagePublisher interface {
	Publish(ctx context.Context, message contract.AmqpMessage) error
}
