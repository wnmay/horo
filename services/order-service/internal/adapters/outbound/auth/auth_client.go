package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// AuthPort defines the interface your app depends on
type AuthPort interface {
	ValidateToken(ctx context.Context, token string) (User, error)
}

// User represents authenticated user info
type User struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

// AuthClient is the concrete implementation of AuthPort
type AuthClient struct {
	BaseURL string
	Client  *http.Client
}

// NewAuthClient creates a new AuthClient with the auth service base URL
func NewAuthClient(baseURL string) AuthPort {
	return &AuthClient{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 5 * time.Second},
	}
}

// ValidateToken calls your friend's auth service to validate the token
func (a *AuthClient) ValidateToken(ctx context.Context, token string) (User, error) {
	if token == "" {
		return User{}, errors.New("token is empty")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", a.BaseURL+"/validate", nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.Client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, errors.New("unauthorized")
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}
