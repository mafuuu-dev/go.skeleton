package request_test

import (
	"backend/core/config"
	"backend/core/constants"
	"backend/core/pkg/logger"
	"backend/core/pkg/request"
	"backend/core/pkg/scope"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"required,gte=18"`
	Email string `json:"email" validate:"required,email"`
}

func setupApp() *fiber.App {
	app := fiber.New()
	return app
}

func setupValidator() *request.Validator {
	cfg := config.Load(constants.ServiceServer)
	sc := &scope.Scope{
		Log: logger.New(cfg),
	}
	return request.NewValidator(sc)
}

func TestValidator_Positive(t *testing.T) {
	app := setupApp()
	v := setupValidator()

	app.Post("/", func(c *fiber.Ctx) error {
		var req TestRequest
		ok, err := v.Validate(c, &req)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, "Alice", req.Name)
		assert.Equal(t, 25, req.Age)
		assert.Equal(t, "alice@example.com", req.Email)
		return c.SendStatus(http.StatusOK)
	})

	body := `{"name": "Alice", "age": 25, "email": "alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestValidator_Negative(t *testing.T) {
	app := setupApp()
	v := setupValidator()

	app.Post("/", func(c *fiber.Ctx) error {
		var req TestRequest
		ok, err := v.Validate(c, &req)
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, c.Response().StatusCode())
		return nil
	})

	body := `{"name": "", "age": 16, "email": "invalid-email"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var result map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.False(t, result["success"].(bool))
	assert.Equal(t, float64(http.StatusUnprocessableEntity), result["code"])
}

func TestValidator_BadJSON(t *testing.T) {
	app := setupApp()
	v := setupValidator()

	app.Post("/", func(c *fiber.Ctx) error {
		req := &TestRequest{}
		ok, err := v.Validate(c, req)
		assert.False(t, ok)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, c.Response().StatusCode())
		return nil
	})

	body := `{"name": "John",}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
