package request

import (
	"backend/core/pkg/scope"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareProvider interface {
	Middleware(handlers ...fiber.Handler) *Handler
}

type HandlerProvider interface {
	Handle() fiber.Handler
}

type Handler struct {
	scope      *scope.Scope
	router     *Builder
	validator  *Validator
	middleware []fiber.Handler
}

func NewHandler(scope *scope.Scope) *Handler {
	return &Handler{
		scope:     scope,
		router:    Router(),
		validator: NewValidator(scope),
	}
}

func (handler *Handler) SC() *scope.Scope {
	return handler.scope
}

func (handler *Handler) Validator() *Validator {
	return handler.validator
}

func (handler *Handler) Instance(h HandlerProvider) []fiber.Handler {
	return handler.router.With(handler.middleware...).Then(h.Handle())
}

func (handler *Handler) Middleware(middleware fiber.Handler) *Handler {
	handler.middleware = append(handler.middleware, middleware)
	return handler
}

func (handler *Handler) Handle() fiber.Handler {
	panic("Handle() must be implemented")
}
