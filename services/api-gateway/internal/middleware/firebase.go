package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	grpcinfra "github.com/wnmay/horo/services/api-gateway/internal/grpc"
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type AuthMiddleware struct {
	authGrpcClient pb.AuthServiceClient
}

func NewAuthMiddleware(clients *grpcinfra.GrpcClients) *AuthMiddleware {
	return &AuthMiddleware{
		authGrpcClient: clients.AuthServiceClient,
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

	// Strip any existing X-User-* headers from incoming request to prevent header injection
	c.Request().Header.Del("X-User-Id")
	c.Request().Header.Del("X-User-Email")
	c.Request().Header.Del("X-User-Roles")

	// Add claims as headers for upstream services
	c.Request().Header.Set("X-User-Id", claimsResp.UserId)
	c.Request().Header.Set("X-User-Email", claimsResp.Email)
	c.Request().Header.Set("X-User-Roles", claimsResp.Role)

	// Continue to next handler (proxy to upstream service)
	return c.Next()
}
