// internal/ports/auth_port.go
package ports

import (
	"context"
)

type Claims struct {
	UserID string
	Email  string
	Role   string
}

type AuthPort interface {
	VerifyIDToken(ctx context.Context, token string) (*Claims, error)
	SetCustomUserClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
}

type AuthService interface {
	GetClaims(ctx context.Context, token string) (*Claims, error)
	SetCustomClaims(ctx context.Context, uid string, customClaims map[string]interface{}) error
}
