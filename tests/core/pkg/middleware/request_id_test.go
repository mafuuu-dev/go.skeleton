package middleware_test

import (
	"backend/core/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestID_GeneratesNewID_WhenMissing(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.RequestID())

	var generatedID string

	app.Get("/", func(c *fiber.Ctx) error {
		id := c.Get(middleware.RequestIDKey)
		assert.NotEmpty(t, id)
		assert.True(t, uuid.Validate(id) == nil, "should be a valid UUID")

		// Проверяем, что доступен через Locals
		assert.Equal(t, id, c.Locals(middleware.RequestIDKey))

		generatedID = id
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, generatedID)
}

func TestRequestID_UsesExistingID(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.RequestID())

	var receivedID string
	existingID := uuid.NewString()

	app.Get("/", func(c *fiber.Ctx) error {
		receivedID = c.Get(middleware.RequestIDKey)
		return c.SendStatus(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(middleware.RequestIDKey, existingID)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, existingID, receivedID)
}
