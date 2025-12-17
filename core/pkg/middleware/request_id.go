package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDKey = "X-Request-Id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(RequestIDKey)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(RequestIDKey, requestID)
		c.Request().Header.Set(RequestIDKey, requestID)

		c.Locals(RequestIDKey, requestID)

		return c.Next()
	}
}
