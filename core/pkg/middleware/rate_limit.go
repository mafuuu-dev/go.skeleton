package middleware

import (
	"backend/core/constants"
	"backend/core/pkg/scope"
	"backend/core/types"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimit(scope *scope.Scope, limit types.RateLimitType) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        limit.Max,
		Expiration: limit.Expiration,
		Next: func(c *fiber.Ctx) bool {
			return scope.Config.Environment.IsDevelopment()
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			userID := c.Locals("user_id")
			if userID != nil {
				return fmt.Sprintf("%v", userID)
			}

			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":  string(constants.TooManyRequestsError),
				"code":   http.StatusTooManyRequests,
				"status": false,
			})
		},
	})
}
