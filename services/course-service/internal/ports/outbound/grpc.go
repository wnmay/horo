package outbound

import (
	"context"

	"github.com/wnmay/horo/services/course-service/internal/domain"
)

type UserProvider interface {
	MapProphetNamesByIDs(ctx context.Context, userIDs []string) ([]domain.ProphetName, error)
	GetProphetName(ctx context.Context, userID string) (string, error)
	GetProphetIDsByNames(ctx context.Context, prophetName string) ([]domain.ProphetName, error)
}
