package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ResponseWrapper() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Capture the response body
		var resBody bytes.Buffer
		c.Response().SetBodyStream(&resBody, len(c.Response().Body()))

		err := c.Next()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Internal Server Error",
			})
		}

		status := c.Response().StatusCode()
		body := resBody.Bytes()

		// Try to parse JSON
		var jsonBody interface{}
		_ = json.Unmarshal(body, &jsonBody)

		// Avoid double wrapping
		if isAlreadyWrapped(jsonBody) {
			return c.Status(status).Send(body)
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

		// Success response
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
