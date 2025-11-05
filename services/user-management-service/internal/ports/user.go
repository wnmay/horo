package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserManagementService interface {
	Register(ctx context.Context, idToken, fullName, role string) error
	GetMe(ctx context.Context, userId string) (*domain.User, error)
}
