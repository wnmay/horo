package jwt

import (
	"fmt"

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
