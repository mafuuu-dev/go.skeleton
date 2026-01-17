package server_usecase

import (
	"backend/core/constants"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct {
	scope *scope.Scope
}

func NewErrorHandler(scope *scope.Scope) *ErrorHandler {
	return &ErrorHandler{scope: scope}
}

func (h *ErrorHandler) Handler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		e := errorsx.Extract(err)

		h.scope.Log.Warnf("Error occurred: %s", e.(*errorsx.Error).ToJSON())

		humanCode, status := e.(*errorsx.Error).GetHuman()
		if humanCode == "" {
			humanCode = string(constants.InternalServerError)

			if status != http.StatusInternalServerError {
				humanCode = e.(*errorsx.Error).GetFirstMessage()
			}
		}

		return response.Error(c, humanCode, status)
	}
}
