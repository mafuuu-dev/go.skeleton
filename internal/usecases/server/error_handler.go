package server_usecase

import (
	"backend/core/constants"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type ErrorHandler struct {
	scope *scope.Scope
}

func NewErrorHandler(scope *scope.Scope) *ErrorHandler {
	return &ErrorHandler{
		scope: scope,
	}
}

func (u *ErrorHandler) Handler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := http.StatusInternalServerError
		message := string(constants.InternalServerError)

		var fe *fiber.Error
		if errors.As(err, &fe) {
			code = fe.Code
			message = fe.Message
		}

		var he *errorsx.HumanError
		if errors.As(err, &he) {
			code = he.Status
			message = string(he.Message)
		}

		u.scope.Log.Warnf("Server error handler: %v", errorsx.EnrichTrace(errorsx.JSONTrace(err), code, err.Error()))

		if code == http.StatusInternalServerError && he == nil {
			message = string(constants.InternalServerError)
		}

		return response.Error(c, message, code)
	}
}
