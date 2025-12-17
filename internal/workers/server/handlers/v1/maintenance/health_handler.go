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

type HealthHandler struct {
	*request.Handler
}

func Health(scope *scope.Scope) []fiber.Handler {
	handler := &HealthHandler{
		Handler: request.NewHandler(scope),
	}

	return handler.
		Middleware(middleware.RateLimit(scope, types.RateLimitType{Max: 20, Expiration: 1 * time.Second})).
		Instance(handler)
}

func (h *HealthHandler) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.Success(c, h.response())
	}
}

func (h *HealthHandler) response() any {
	return struct {
		Health string `json:"health"`
	}{
		Health: "ok",
	}
}
