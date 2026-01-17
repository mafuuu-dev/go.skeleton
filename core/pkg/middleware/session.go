package middleware

import (
	"backend/core/constants"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/scope"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PayloadProvider interface {
	IsCanAction(c *fiber.Ctx) error
	SavePayloadToLocal(c *fiber.Ctx) error
}

func Session(scope *scope.Scope, provider PayloadProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := provider.SavePayloadToLocal(c); err != nil {
			scope.Log.Warn(
				errorsx.WrapHuman(err, constants.InternalServerError, http.StatusInternalServerError).(*errorsx.Error).
					ToJSON(),
			)

			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error":  string(constants.InternalServerError),
				"code":   http.StatusInternalServerError,
				"status": false,
			})
		}

		if err := provider.IsCanAction(c); err != nil {
			scope.Log.Warn(
				errorsx.WrapHuman(err, constants.ForbiddenError, http.StatusForbidden).(*errorsx.Error).ToJSON(),
			)

			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error":  string(constants.ForbiddenError),
				"code":   http.StatusForbidden,
				"status": false,
			})
		}

		return c.Next()
	}
}
