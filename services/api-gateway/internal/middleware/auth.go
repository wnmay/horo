package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader || token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	url := fmt.Sprintf("%s/api/auth/verify-token?token=%s", a.authServiceAddr, token)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("[AuthMiddleware] error calling auth service: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to contact auth service",
			"details": err.Error(),
		})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("[AuthMiddleware] auth service returned %d: %s\n", res.StatusCode, string(bodyBytes))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to read auth service response",
			"details": err.Error(),
		})
	}

	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "failed to unmarshal auth response",
			"details": err.Error(),
		})
	}

	// Validate contents of the auth response
	if authResponse.UserID == "" || authResponse.Email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "failed to authenticate",
			"details": string(body),
		})
	}

	// Set claims for downstream handlers
	c.Request().Header.Set("X-User-Id", authResponse.UserID)
	c.Request().Header.Set("X-User-Email", authResponse.Email)
	c.Request().Header.Set("X-User-Role", authResponse.Role)

	c.Locals("userId", authResponse.UserID)
	c.Locals("userEmail", authResponse.Email)
	c.Locals("userRole", authResponse.Role)

	log.Printf("[AuthMiddleware] Auth OK: %+v\n", authResponse)

	return c.Next()
}
