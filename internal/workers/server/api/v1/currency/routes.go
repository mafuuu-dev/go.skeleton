package v1_currency

import (
	"backend/core/pkg/scope"
	"backend/internal/workers/server/handlers/v1/currency"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(router fiber.Router, scope *scope.Scope) {
	player := router.Group("/currency")

	player.Get("/", v1_handler_currency.GetCurrency(scope)...)
}
