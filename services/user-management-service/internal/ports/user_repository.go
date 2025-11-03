package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserRepositoryPort interface {
	Save(ctx context.Context, user domain.User) error
}
