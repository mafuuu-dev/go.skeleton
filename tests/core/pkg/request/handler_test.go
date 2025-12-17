package request_test

import (
	"backend/core/pkg/request"
	"backend/core/pkg/scope"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	fn fiber.Handler
}

func (m mockHandler) Handle() fiber.Handler {
	return m.fn
}

func TestNewHandler_Initialization(t *testing.T) {
	sc := &scope.Scope{}
	h := request.NewHandler(sc)

	assert.NotNil(t, h)
	assert.Equal(t, sc, h.SC())
	assert.NotNil(t, h.Validator())
}

func TestHandler_Middleware_ChainApplied(t *testing.T) {
	app := fiber.New()
	sc := &scope.Scope{}
	h := request.NewHandler(sc)

	var order []string

	h.Middleware(func(c *fiber.Ctx) error {
		order = append(order, "m1")
		return c.Next()
	})
	h.Middleware(func(c *fiber.Ctx) error {
		order = append(order, "m2")
		return c.Next()
	})

	mock := mockHandler{fn: func(c *fiber.Ctx) error {
		order = append(order, "handler")
		return c.SendStatus(http.StatusOK)
	}}

	app.Get("/", h.Instance(mock)...)
	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []string{"m1", "m2", "handler"}, order)
}

func TestHandler_Instance_BuildsHandlerChain(t *testing.T) {
	app := fiber.New()
	sc := &scope.Scope{}
	h := request.NewHandler(sc)

	called := false
	mock := mockHandler{fn: func(c *fiber.Ctx) error {
		called = true
		return c.SendStatus(http.StatusOK)
	}}

	h.Middleware(func(c *fiber.Ctx) error {
		c.Set("X-Test", "ok")
		return c.Next()
	})

	handlers := h.Instance(mock)
	assert.Len(t, handlers, 2)

	app.Get("/", handlers...)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, called)
	assert.Equal(t, "ok", resp.Header.Get("X-Test"))
}

func TestHandler_Handle_Panic(t *testing.T) {
	h := &request.Handler{}
	assert.Panics(t, func() {
		h.Handle()
	})
}
