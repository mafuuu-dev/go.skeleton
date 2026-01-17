package response

import (
	"backend/core/constants"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/scope"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
}

func Error(c *fiber.Ctx, err string, code int) error {
	return c.Status(code).JSON(ErrorResponse{
		Error:   err,
		Code:    code,
		Success: false,
	})
}

func NotValidRequest(c *fiber.Ctx, scope *scope.Scope, err error) error {
	scope.Log.Warn(errorsx.WrapJSON(err, "Request validation error"))
	return Error(c, string(constants.UnprocessableEntityError), http.StatusUnprocessableEntity)
}
