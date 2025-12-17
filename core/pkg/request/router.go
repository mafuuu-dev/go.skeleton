package request

import (
	"github.com/gofiber/fiber/v2"
)

type Builder struct {
	handlers []fiber.Handler
}

func Router() *Builder {
	return &Builder{
		handlers: make([]fiber.Handler, 0, 4),
	}
}

func (builder *Builder) With(handlers ...fiber.Handler) *Builder {
	builder.handlers = append(builder.handlers, handlers...)
	return builder
}

func (builder *Builder) Then(handler fiber.Handler) []fiber.Handler {
	return append(builder.handlers, handler)
}
