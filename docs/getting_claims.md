# Accessing User Claims from Requests

This document explains how upstream services can retrieve authenticated user information from HTTP headers injected by the API Gateway or Auth middleware.

## Overview

When a client sends a request with a valid Firebase ID token in the `Authorization` header, the API Gateway validates the token and extracts the user claims.  
After verification, it forwards the request to the upstream service **with trusted headers** that contain user identity information.

## Injected Headers

After successful authentication, the gateway adds the following headers:

| Header         | Description                                       | Example            |
| -------------- | ------------------------------------------------- | ------------------ |
| `X-User-Id`   | Unique Firebase user ID                           | `6nFZV...W2bA`     |
| `X-User-Email` | User’s email (if available)                       | `user@example.com` |
| `X-User-Name`  | Display name (optional)                           | `Jane Doe`         |
| `X-User-Role`  | Comma-separated list of user roles or permissions | `customer`         |

> Note: The gateway strips these headers from all incoming requests before injection to prevent header spoofing.

## How to Use in Your Service

Example (Go/fiber):

```go
app.Get("/profile", func(c *fiber.Ctx) error {
		// Read headers injected by the gateway
		userID := c.Get("X-User-Id")
		email := c.Get("X-User-Email")
		role := c.Get("X-User-Role")
        ...
})
```
