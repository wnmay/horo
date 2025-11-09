package http_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userManagementURL string
	httpClient        *http.Client
	validator         *validator.Validate
}

type RegisterRequest struct {
	IdToken  string `json:"idToken" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type RegisterHTTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func NewUserHandler() *UserHandler {
	validate := validator.New()
	userManagementURL := os.Getenv("USER_MANAGEMENT_SERVICE_URL")
	if userManagementURL == "" {
		userManagementURL = "http://localhost:3003"
	}

	return &UserHandler{
		userManagementURL: userManagementURL,
		httpClient:        &http.Client{},
		validator:         validate,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
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

	jsonData, err := json.Marshal(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to marshal request",
		})
	}

	// Make HTTP POST request
	url := fmt.Sprintf("%s/users/register", h.userManagementURL)
	resp, err := h.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to register user",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response",
		})
	}

	// Parse response
	var registerResp RegisterHTTPResponse
	if err := json.Unmarshal(body, &registerResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse response",
		})
	}

	// Handle non-success status codes
	if resp.StatusCode != http.StatusCreated {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "failed to register user",
			"details": registerResp.Message,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
	})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	// Forward the request to the user management service
	url := fmt.Sprintf("%s/users/me", h.userManagementURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to create request",
			"details": err.Error(),
		})
	}

	// Set the Authorization header
	req.Header.Set("Authorization", authHeader)

	// Perform the HTTP request
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to fetch user info",
			"details": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response",
		})
	}

	// Handle non-success status codes
	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"error":   "failed to fetch user info",
			"details": string(body),
		})
	}

	// Return the response body as JSON
	var userData map[string]interface{}
	if err := json.Unmarshal(body, &userData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to parse user info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(userData)
}
