// internal/app/user_management_service.go
package app

import (
	"context"
	"fmt"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
)

type UserManagementService struct {
	auth ports.AuthPort
	repo ports.UserRepositoryPort
}

func NewUserManagementService(auth ports.AuthPort, repo ports.UserRepositoryPort) *UserManagementService {
	return &UserManagementService{auth: auth, repo: repo}
}

func (s *UserManagementService) Register(ctx context.Context, idToken, fullName, role string) error {
	claims, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return fmt.Errorf("invalid firebase token: %w", err)
	}

	user := domain.User{
		ID:       claims.UserID,
		FullName: fullName,
		Email:    claims.Email,
		Role:     role,
	}

	return s.repo.Save(ctx, user)
}
