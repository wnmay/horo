package ports

import (
	"context"
)

type UserManagementService interface {
	Register(ctx context.Context, idToken, fullName, role string) error
}
