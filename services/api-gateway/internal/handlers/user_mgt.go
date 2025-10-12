package handlers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	grpcinfra "github.com/wnmay/horo/services/api-gateway/internal/grpc"
	pb "github.com/wnmay/horo/shared/proto/user-management"
)

type UserHandler struct {
	UserManagementClient pb.UserManagementServiceClient
	validator            *validator.Validate
}

type RegisterRequest struct {
	FullName string `json:"fullName" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

func NewUserHandler(client *grpcinfra.GrpcClients) *UserHandler {
	return &UserHandler{
		UserManagementClient: client.UserManagementClient,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// Extract Bearer token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	// Remove "Bearer " prefix
	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	if idToken == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	// Parse request body
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Call gRPC client
	ctx := c.Context()

	grpcReq := &pb.RegisterRequest{
		FirebaseToken: idToken,
		FullName:      req.FullName,
		Role:          req.Role,
	}

	_, err := h.UserManagementClient.Register(ctx, grpcReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to register user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}
