package outbound_port

import (
	"context"

	"github.com/wnmay/horo/services/chat-service/internal/domain"
)

type UserProvider interface {
	MapUserNamesByIDs(ctx context.Context, userIDs []string) (map[string]*domain.User, error)
}
