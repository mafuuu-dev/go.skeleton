package response

import (
	"backend/core/pkg/response"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/success", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{"hello": "world"})
	})

	req := httptest.NewRequest("GET", "/success", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, true, body["success"])
	assert.Equal(t, "world", body["data"].(map[string]interface{})["hello"])
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()

	app.Get("/error", func(c *fiber.Ctx) error {
		return response.Error(c, "something went wrong", 400)
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, false, body["success"])
	assert.Equal(t, "something went wrong", body["error"])
}
