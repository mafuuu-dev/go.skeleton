package api

import (
	"backend/core/pkg/scope"
	"backend/internal/workers/server/api/v1/currency"
	"backend/internal/workers/server/api/v1/maintenance"
	"backend/internal/workers/server/handlers"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, scope *scope.Scope) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1_maintenance.RegisterRoutes(v1, scope)
	v1_currency.RegisterRoutes(v1, scope)

	if scope.Config.Environment.IsDevelopment() {
		api.Get("/", handler.RouteList(scope, app)...)
	}
}
