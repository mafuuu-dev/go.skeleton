package v1_maintenance

import (
	"backend/core/pkg/scope"
	"backend/internal/workers/server/handlers/v1/maintenance"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, scope *scope.Scope) {
	group := router.Group("/maintenance")

	group.Get("/ping", v1_handler_maintenance.Ping(scope)...)
	group.Get("/health", v1_handler_maintenance.Health(scope)...)
}
