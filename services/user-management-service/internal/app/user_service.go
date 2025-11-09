// internal/app/user_management_service.go
package app

import (
	"context"
	"fmt"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserManagementService struct {
	authClient  ports.AuthPort
	repo        ports.UserRepositoryPort
	prophetRepo ports.ProphetRepoPort
}

func NewUserManagementService(authClient ports.AuthPort, repo ports.UserRepositoryPort, prophetRepo ports.ProphetRepoPort) *UserManagementService {
	return &UserManagementService{authClient: authClient, repo: repo, prophetRepo: prophetRepo}
}

func (s *UserManagementService) Register(ctx context.Context, idToken, fullName, role string) error {
	claims, err := s.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return fmt.Errorf("invalid firebase token: %w", err)
	}

	uid := claims.UserID
	customClaims := map[string]interface{}{
		"role": role,
	}

	if err := s.authClient.SetCustomUserClaims(ctx, uid, customClaims); err != nil {
		return status.Errorf(codes.Internal, "failed to set custom claims: %v", err)
	}

	user := domain.User{
		ID:       claims.UserID,
		FullName: fullName,
		Email:    claims.Email,
		Role:     role,
	}

	if user.Role == "customer" {
		return s.repo.Save(ctx, user)
	}

	prophet := domain.Prophet{
		User:    &user,
		Balance: 0,
	}
	return s.prophetRepo.Save(ctx, prophet)
}

func (s *UserManagementService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
