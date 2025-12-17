package handler

import (
	"backend/core/pkg/request"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"backend/internal/usecases/server"

	"github.com/gofiber/fiber/v2"
)

type RouteListHandler struct {
	*request.Handler
	app *fiber.App
}

func RouteList(scope *scope.Scope, app *fiber.App) []fiber.Handler {
	handler := &RouteListHandler{
		Handler: request.NewHandler(scope),
		app:     app,
	}

	return handler.Instance(handler)
}

func (h *RouteListHandler) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.Success(c, server_usecase.NewRouteList(h.app).Get())
	}
}
