// internal/app/service.go
package app

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/wnmay/horo/services/user-management-service/internal/config"
	"github.com/wnmay/horo/services/user-management-service/internal/domain"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
)

type userManagementService struct {
	repo               ports.UserRepository
	firebaseAuthClient *auth.Client
}

func NewUserManagementService(ctx context.Context, repo ports.UserRepository, cfg *config.Config) ports.UserManagementService {
	firebaseAuthClient := InitFirebase(ctx, cfg)
	return &userManagementService{
		repo:               repo,
		firebaseAuthClient: firebaseAuthClient,
	}
}

func (s *userManagementService) Register(ctx context.Context, idToken, fullName, role string) error {
	// Verify the Firebase ID token
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return fmt.Errorf("invalid firebase token: %w", err)
	}

	// Extract user info from token
	firebaseUID := token.UID
	email := token.Claims["email"].(string)

	user := domain.User{
		ID:       firebaseUID, // Use Firebase UID as primary key
		FullName: fullName,
		Email:    email,
		Role:     role,
	}

	err = s.repo.Save(user)
	return err
}

func (s *userManagementService) ValidateFirebaseToken(firebaseToken string) (*domain.Claims, error) {
	ctx := context.Background()
	token, err := s.firebaseAuthClient.VerifyIDToken(ctx, firebaseToken)
	if err != nil {
		return nil, fmt.Errorf("invalid Firebase token: %w", err)
	}

	// Extract claims safely
	userID := token.UID
	email, _ := token.Claims["email"].(string)
	role, _ := token.Claims["role"].(string) // Custom claim

	return &domain.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
	}, nil
}
