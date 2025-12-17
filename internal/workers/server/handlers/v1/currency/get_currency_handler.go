package v1_handler_currency

import (
	"backend/core/pkg/errorsx"
	"backend/core/pkg/middleware"
	"backend/core/pkg/request"
	"backend/core/pkg/response"
	"backend/core/pkg/scope"
	"backend/core/types"
	"backend/internal/constants"
	"backend/internal/domain/currency/enum"
	"backend/internal/usecases/currency"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type getCurrencyHandler struct {
	Code string `json:"code" validate:"required"`
}

type GetCurrencyHandler struct {
	*request.Handler
}

func GetCurrency(scope *scope.Scope) []fiber.Handler {
	handler := &GetCurrencyHandler{
		Handler: request.NewHandler(scope),
	}

	return handler.
		Middleware(middleware.RateLimit(scope, types.RateLimitType{Max: 30, Expiration: 1 * time.Minute})).
		Instance(handler)
}

func (h *GetCurrencyHandler) Handle() fiber.Handler {
	return func(c *fiber.Ctx) error {
		model := &getCurrencyHandler{}
		if ok, err := h.Validator().Validate(c, model); !ok {
			return errorsx.Error(err)
		}

		currency, err := currency_usecase.NewGetCurrencyByCode(h.SC()).Get(currency_enum.Code(model.Code))
		if err != nil {
			return errorsx.Humanf(err, internal_constants.GetCurrencyByCodeError, http.StatusBadRequest)
		}

		return response.Success(c, currency)
	}
}
