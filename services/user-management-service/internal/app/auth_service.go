// internal/app/auth_app_service.go
package app

import (
	"context"

	"github.com/wnmay/horo/services/user-management-service/internal/ports"
)

type AuthService struct {
	authPort ports.AuthPort
}

func NewAuthService(authPort ports.AuthPort) *AuthService {
	return &AuthService{authPort: authPort}
}

func (s *AuthService) GetClaims(ctx context.Context, token string) (*ports.Claims, error) {
	claims, err := s.authPort.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *AuthService) SetCustomClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error {
	return s.authPort.SetCustomUserClaims(ctx, uid, customClaims)
}
