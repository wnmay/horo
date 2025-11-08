package jwt

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Secret string
	Issuer string
}

var config Config

func Init(cfg Config) {
	if cfg.Secret == "" {
		panic("jwt: Secret cannot be empty")
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "default-issuer"
	}
	config = cfg
}

type JWTClaims struct {
	UserID     string `json:"user_id"`
	ProphetID  string `json:"prophet_id"`
	CustomerID string `json:"customer_id"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if claims.Issuer != config.Issuer {
			return "", fmt.Errorf("invalid issuer")
		}
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token")
}

func ExtractClaims(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// VerifyFirebaseToken verifies a Firebase ID token and returns the user ID
func VerifyFirebaseToken(ctx context.Context, token string) (string, error) {
	// Get Firebase auth client from context or initialize it
	// This assumes you've initialized Firebase in your main.go
	client, err := getFirebaseAuthClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get Firebase auth client: %w", err)
	}

	// Verify the ID token
	t, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("failed to verify Firebase token: %w", err)
	}

	return t.UID, nil
}

// Global Firebase auth client
var firebaseAuthClient *auth.Client

// InitFirebase initializes the Firebase auth client
func InitFirebase(client *auth.Client) {
	firebaseAuthClient = client
}

// getFirebaseAuthClient returns the Firebase auth client
func getFirebaseAuthClient(ctx context.Context) (*auth.Client, error) {
	if firebaseAuthClient == nil {
		return nil, fmt.Errorf("Firebase auth client not initialized. Call jwt.InitFirebase() first")
	}
	return firebaseAuthClient, nil
}
