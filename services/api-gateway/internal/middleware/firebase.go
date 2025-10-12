package middleware

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	grpcinfra "github.com/wnmay/horo/services/api-gateway/internal/grpc"
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type AuthMiddleware struct {
	authGrpcClient pb.UserManagementServiceClient
}

func NewAuthMiddleware(clients *grpcinfra.GrpcClients) *AuthMiddleware {
	return &AuthMiddleware{
		authGrpcClient: clients.UserManagementClient,
	}
}

func (a *AuthMiddleware) AddClaims(c *fiber.Ctx) error {
	// Extract Bearer token
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	// Call gRPC to validate token and get claims
	ctx := c.Context()
	grpcReq := &pb.GetClaimsRequest{
		FirebaseToken: token,
	}

	claimsResp, err := a.authGrpcClient.GetClaims(ctx, grpcReq)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid or expired token",
			"details": err.Error(),
		})
	}

	// Parse existing body
	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		// If body is empty or invalid, create new map
		body = make(map[string]interface{})
	}

	// Add claims to body
	body["claims"] = map[string]string{
		"user_id": claimsResp.UserId,
		"email":   claimsResp.Email,
		"role":    claimsResp.Role,
	}

	// Update request body with claims
	c.Request().SetBody([]byte(mustMarshal(body)))
	c.Request().Header.SetContentType("application/json")

	// Continue to next handler
	return c.Next()
}

func mustMarshal(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
