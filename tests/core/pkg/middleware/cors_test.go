package middleware_test

import (
	"backend/core/config"
	"backend/core/pkg/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupCors(corsOrigins string) *fiber.App {
	cfg := &config.Config{
		CORSOrigins: corsOrigins,
	}

	app := fiber.New()
	app.Use(middleware.Cors(cfg))
	app.All("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	return app
}

func TestCorsMiddleware_AllowsConfiguredOrigin(t *testing.T) {
	app := setupCors("http://localhost:3000")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
}

func TestCorsMiddleware_BlocksOtherOrigin(t *testing.T) {
	app := setupCors("http://example.com")

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://notallowed.com")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, "", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCorsMiddleware_HandlesOptionsPreflight(t *testing.T) {
	app := setupCors("http://localhost:3000")

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "http://localhost:3000", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
}
