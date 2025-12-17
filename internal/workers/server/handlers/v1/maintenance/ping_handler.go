package v1_handler_maintenance

import (
	"backend/core/pkg/middleware"
	"backend/core/pkg/request"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"backend/core/types"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PingHandler struct {
	*request.Handler
}

func Ping(scope *scope.Scope) []fiber.Handler {
	handler := &PingHandler{
		Handler: request.NewHandler(scope),
	}

	return handler.
		Middleware(middleware.RateLimit(scope, types.RateLimitType{Max: 200, Expiration: 1 * time.Second})).
		Instance(handler)
}

func (h *PingHandler) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		return response.Success(c, fiber.Map{})
	}
}
