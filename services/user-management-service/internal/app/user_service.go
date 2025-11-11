// internal/app/user_management_service.go
package app

import (
	"context"
	"fmt"
	"log"

	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserManagementService struct {
	authClient ports.AuthPort
	repo       ports.UserRepositoryPort
}

func NewUserManagementService(authClient ports.AuthPort, repo ports.UserRepositoryPort) *UserManagementService {
	return &UserManagementService{authClient: authClient, repo: repo}
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

	log.Printf("Setting custom claims for user %s: %v", uid, customClaims)
	if err := s.authClient.SetCustomUserClaims(ctx, uid, customClaims); err != nil {
		return status.Errorf(codes.Internal, "failed to set custom claims: %v", err)
	}

	user := domain.User{
		ID:       claims.UserID,
		FullName: fullName,
		Email:    claims.Email,
		Role:     role,
	}
	return s.repo.Save(ctx, user)
}

func (s *UserManagementService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := s.repo.FindById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserManagementService) UpdateFullName(ctx context.Context, userID string, newUsername string) (*domain.User, error) {
	update := map[string]interface{}{"fullname": newUsername}
	return s.repo.Update(ctx, userID, update)
}

func (s *UserManagementService) GetProphetNames(ctx context.Context, userIDs []string) ([]*domain.ProphetName, error) {
	prophetNames, err := s.repo.FindProphetNames(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return prophetNames, nil
}

func (s *UserManagementService) GetProphetName(ctx context.Context, userID string) (string, error) {
	prophet, err := s.repo.FindById(ctx, userID)
	if err != nil {
		return "", err
	}
	return prophet.FullName, nil
}

func (s *UserManagementService) SearchProphetIdsByName(ctx context.Context, prophetName string) ([]*domain.ProphetName, error) {
	prophetIds, err := s.repo.SearchProphetIdsByName(ctx, prophetName)
	if err != nil {
		return nil, err
	}
	return prophetIds, nil
}

func (s *UserManagementService) MapUserNames(ctx context.Context, userIDs []string) ([]*domain.UserName, error) {
	userNames, err := s.repo.MapUserNames(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return userNames, nil
}
