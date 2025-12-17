package middleware_test

import (
	"backend/core/config"
	"backend/core/constants"
	"backend/core/pkg/logger"
	"backend/core/pkg/middleware"
	"backend/core/pkg/scope"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type mockPayloadProvider struct {
	saveErr error
	canErr  error
}

func (m *mockPayloadProvider) SavePayloadToLocal(c *fiber.Ctx) error {
	return m.saveErr
}

func (m *mockPayloadProvider) IsCanAction(c *fiber.Ctx) error {
	return m.canErr
}

func setupSession(provider middleware.PayloadProvider) *fiber.App {
	cfg := config.Load(constants.ServiceServer)
	sc := &scope.Scope{
		Log: logger.New(cfg),
	}
	app := fiber.New()
	app.Use(middleware.Session(sc, provider))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	return app
}

func TestSession_Success(t *testing.T) {
	provider := &mockPayloadProvider{}
	app := setupSession(provider)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSession_SavePayloadError(t *testing.T) {
	provider := &mockPayloadProvider{
		saveErr: errors.New("failed to save payload"),
	}
	app := setupSession(provider)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestSession_IsCanActionError(t *testing.T) {
	provider := &mockPayloadProvider{
		canErr: errors.New("not allowed"),
	}
	app := setupSession(provider)

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
