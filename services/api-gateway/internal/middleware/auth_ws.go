package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (a *AuthMiddleware) AddClaimsWS(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token query parameter",
		})
	}

	// Call your auth service to verify token (same as AddClaims)
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
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "failed to unmarshal auth response",
			"details": err.Error(),
		})
	}

	// Attach user info to context
	c.Locals("userId", authResponse.UserID)
	c.Locals("email", authResponse.Email)
	c.Locals("role", authResponse.Role)

	return c.Next()
}
