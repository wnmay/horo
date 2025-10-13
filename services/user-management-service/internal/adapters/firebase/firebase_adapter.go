// internal/adapters/firebase/auth_adapter.go
package firebase

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
)

type AuthAdapter struct {
	client *auth.Client
}

func NewAuthAdapter(client *auth.Client) *AuthAdapter {
	return &AuthAdapter{client: client}
}

func (f *AuthAdapter) VerifyIDToken(ctx context.Context, token string) (*ports.Claims, error) {
	t, err := f.client.VerifyIDToken(ctx, token)
	if err != nil {
		return nil, err
	}

	claims := &ports.Claims{
		UserID: t.UID,
	}

	// extract email (Firebase standard claim)
	if email, ok := t.Claims["email"].(string); ok {
		claims.Email = email
	}

	// extract custom "role" claim if set in Firebase
	if role, ok := t.Claims["role"].(string); ok {
		claims.Role = role
	}

	return claims, nil
}
