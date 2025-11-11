package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserManagementService interface {
	Register(ctx context.Context, idToken, fullName, role string) error
	GetMe(ctx context.Context, userID string) (*domain.User, error)
	UpdateFullName(ctx context.Context, userID string, newUsername string) (*domain.User, error)
	GetProphetNames(ctx context.Context, userIDs []string) ([]*domain.ProphetName, error)
	GetProphetName(ctx context.Context, userID string) (string, error)
	SearchProphetIdsByName(ctx context.Context, prophetName string) ([]*domain.ProphetName, error)
	MapUserNames(ctx context.Context, userIDs []string) ([]*domain.UserName, error)
}
