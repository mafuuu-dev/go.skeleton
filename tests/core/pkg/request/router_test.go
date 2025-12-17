package request_test

import (
	"backend/core/pkg/request"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRouter_BasicInitialization(t *testing.T) {
	b := request.Router()
	assert.NotNil(t, b, "Router() should not return nil")
}

func TestBuilder_With_And_Then(t *testing.T) {
	app := fiber.New()
	builder := request.Router()

	var order []string

	m1 := func(c *fiber.Ctx) error {
		order = append(order, "m1")
		return c.Next()
	}
	m2 := func(c *fiber.Ctx) error {
		order = append(order, "m2")
		return c.Next()
	}
	final := func(c *fiber.Ctx) error {
		order = append(order, "final")
		return c.SendStatus(http.StatusOK)
	}

	handlers := builder.With(m1, m2).Then(final)
	app.Get("/", handlers...)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, []string{"m1", "m2", "final"}, order)
}

func TestBuilder_Then_WithoutMiddleware(t *testing.T) {
	app := fiber.New()
	builder := request.Router()

	called := false
	final := func(c *fiber.Ctx) error {
		called = true
		return c.SendStatus(http.StatusOK)
	}

	app.Get("/", builder.Then(final)...)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, called)
}
