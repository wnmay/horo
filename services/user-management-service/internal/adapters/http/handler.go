package http

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/wnmay/horo/services/user-management-service/internal/ports"
)

type HTTPHandler struct {
	userService ports.UserManagementService
	authService ports.AuthService
}

func NewHTTPHandler(userService ports.UserManagementService, authService ports.AuthService) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		authService: authService,
	}
}

type RegisterRequest struct {
	IdToken  string `json:"idToken" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func (h *HTTPHandler) SetupRoutes(app *fiber.App) {
	// Apply middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "user-management-service",
		})
	})

	// User routes
	users := app.Group("/users")
	users.Post("/register", h.Register)
	users.Get("/me", h.GetMe)
}

func (h *HTTPHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Validate required fields
	if req.IdToken == "" || req.FullName == "" || req.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(RegisterResponse{
			Success: false,
			Message: "Missing required fields: idToken, fullName, and role are required",
		})
	}

	ctx := context.Background()
	err := h.userService.Register(ctx, req.IdToken, req.FullName, req.Role)
	if err != nil {
		log.Printf("Registration failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(RegisterResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

func StartHTTPServer(handler *HTTPHandler, port string) error {
	app := fiber.New(fiber.Config{
		AppName: "User Management Service",
	})

	handler.SetupRoutes(app)

	log.Printf("HTTP server starting on port %s", port)
	return app.Listen(":" + port)
}

func (h *HTTPHandler) GetMe(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	// Extract token from header
	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}
	token := authHeader[len(prefix):]

	// Verify and extract claims
	claims, err := h.authService.GetClaims(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "invalid or expired token",
			"details": err.Error(),
		})
	}

	// Return claims as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"uid":   claims.UserID,
		"email": claims.Email,
		"role":  claims.Role,
	})
}
