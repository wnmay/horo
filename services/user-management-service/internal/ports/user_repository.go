package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserRepositoryPort interface {
	Save(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, userId string) (*domain.User, error)
	Update(ctx context.Context, userID string, update map[string]interface{}) (*domain.User, error)
}
