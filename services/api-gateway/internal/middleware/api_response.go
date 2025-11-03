package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ResponseWrapper() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Call the next handler first
		err := c.Next()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Internal Server Error",
			})
		}

		// Get the response after handler execution
		status := c.Response().StatusCode()
		body := c.Response().Body()

		// If no body, return as is
		if len(body) == 0 {
			return nil
		}

		// Try to parse JSON
		var jsonBody interface{}
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			// Not JSON, return as is
			return nil
		}

		// Avoid double wrapping
		if isAlreadyWrapped(jsonBody) {
			return nil
		}

		if status >= 400 {
			// Error response
			var errValue interface{}
			if m, ok := jsonBody.(map[string]interface{}); ok && len(m) == 1 && m["error"] != nil {
				errValue = m["error"]
			} else if jsonBody != nil {
				errValue = jsonBody
			} else {
				errValue = string(body)
			}

			resp := fiber.Map{
				"error":   errValue,
				"message": http.StatusText(status),
			}
			return c.Status(status).JSON(resp)
		}

		// Success response - wrap the data
		var dataValue interface{}
		switch v := jsonBody.(type) {
		case nil:
			dataValue = []interface{}{} // empty array if no data
		default:
			dataValue = v
		}

		resp := fiber.Map{
			"data":    dataValue,
			"message": "OK",
		}
		return c.Status(status).JSON(resp)
	}
}

func isAlreadyWrapped(b interface{}) bool {
	m, ok := b.(map[string]interface{})
	if !ok {
		return false
	}
	_, hasData := m["data"]
	_, hasError := m["error"]
	return hasData || hasError
}
