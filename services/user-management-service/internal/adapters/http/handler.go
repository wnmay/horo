package http

import (
	"context"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
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

type VerifyTokenResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
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
	api := app.Group("/api")

	// User routes
	users := api.Group("/users")
	users.Post("/register", h.Register)
	users.Patch("/update-name", h.UpdateUsername)
	auth := api.Group("/auth")
	auth.Get("/verify-token", h.VerifyToken)
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

func (h *HTTPHandler) VerifyToken(c *fiber.Ctx) error {
	token := c.Query("token")
	log.Println("token", token)
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "token is required",
		})
	}

	claims, err := h.authService.GetClaims(c.Context(), token)
	if err != nil {
		log.Println("error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	return c.Status(fiber.StatusOK).JSON(VerifyTokenResponse{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
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

type updateUsernameReq struct {
	FullName string `json:"fullname" validate:"required,min=2,max=80"`
}

type updateUsernameResponse struct {
	Success     bool   `json:"success"`
	NewUsername string `json:"new_username,omitempty"`
	Message     string `json:"message,omitempty"`
}

func (h *HTTPHandler) UpdateUsername(c *fiber.Ctx) error {
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

	uid := claims.UserID

	// parse body
	var req updateUsernameReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json body")
	}
	req.FullName = strings.TrimSpace(req.FullName)

	// validate
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "fullname is required (2-80 chars)")
	}

	// call service
	u, err := h.userService.UpdateFullName(c.Context(), uid, req.FullName)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "unable to update name")
	}
	if u == nil {
		log.Println(uid)
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	// response
	return c.Status(fiber.StatusOK).JSON(updateUsernameResponse{
		Success:     true,
		NewUsername: u.FullName,
		Message:     "Username updated successfully.",
	})
}
