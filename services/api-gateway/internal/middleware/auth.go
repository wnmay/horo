package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authServiceAddr string
}

type AuthResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func NewAuthMiddleware(authServiceAddr string) *AuthMiddleware {
	return &AuthMiddleware{
		authServiceAddr: authServiceAddr,
	}
}

func (a *AuthMiddleware) AddClaims(c *fiber.Ctx) error {
	// Extract Bearer token
	authHeader := c.Get("Authorization")
	var token string

	if authHeader != "" {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header format",
			})
		}
	}

	// extract from query param for connecting ws
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header or token query param",
		})
	}

	// Call gRPC to validate token and get claims
	res, err := http.Get(fmt.Sprintf("%s/api/auth/verify-token?token=%s", a.authServiceAddr, token))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid or expired token",
			"details": err.Error(),
		})
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid or expired token",
			"details": err.Error(),
		})
	}
	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "failed to unmarshal auth response",
			"details": err.Error(),
		})
	}

	// Continue to next handler (proxy to upstream service)
	return c.Next()
}
