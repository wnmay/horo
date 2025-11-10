package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserRepositoryPort interface {
	Save(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, userID string) (*domain.User, error)
	FindProphetNames(ctx context.Context, userIDs []string) ([]*domain.ProphetName, error)
	SearchProphetIdsByName(ctx context.Context, prophetName string) ([]*domain.ProphetName, error)
	MapUserNames(ctx context.Context, userIDs []string) ([]*domain.UserName, error)
}
