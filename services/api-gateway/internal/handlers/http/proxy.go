package http_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ProxyRequest(c *fiber.Ctx, client *http.Client, method, serviceURL, path string) error {
	// Build target URL
	targetURL := serviceURL + path

	// Get query parameters
	query := c.Request().URI().QueryString()
	if len(query) > 0 {
		targetURL += "?" + string(query)
	}

	// Create request body
	var body io.Reader
	if method == "POST" || method == "PUT" || method == "PATCH" {
		body = bytes.NewReader(c.Body())
	}

	// Create HTTP request
	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create request",
		})
	}

	// Copy headers (including Authorization for auth)
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		errstr := fmt.Sprintf("Failed to reach payment service: %v\n", err) // ADD THIS
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": errstr,
		})
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read response",
		})
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Return response
	c.Status(resp.StatusCode)

	// Try to parse as JSON, if fails return as raw
	var jsonResp interface{}
	if err := json.Unmarshal(respBody, &jsonResp); err == nil {
		return c.JSON(jsonResp)
	}

	return c.Send(respBody)
}
