package ports

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
)

type UserManagementService interface {
	ValidateFirebaseToken(firebaseToken string) (*domain.Claims, error)
	Register(ctx context.Context, idToken, fullName, role string) error
}
