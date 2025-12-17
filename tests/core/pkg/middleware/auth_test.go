package middleware_test

import (
	"backend/core/config"
	"backend/core/pkg/middleware"
	"backend/core/pkg/scope"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const TEST_SECRET string = "test_secret"

type AuthProvider func(c *fiber.Ctx) bool

func (f AuthProvider) IsValidToken(c *fiber.Ctx) bool {
	return f(c)
}

func generateToken(secret string, userID int64, accountID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func setupAuth(secret string) *fiber.App {
	app := fiber.New()

	sc := &scope.Scope{
		Config: &config.Config{
			JWTSecret: secret,
		},
	}

	app.Get(
		"/protected",
		middleware.Auth(sc, AuthProvider(func(c *fiber.Ctx) bool {
			return true
		})),
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "ok"})
		},
	)

	return app
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	app := setupAuth(TEST_SECRET)

	token, err := generateToken(TEST_SECRET, 123, 123)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	app := setupAuth(TEST_SECRET)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.value")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	app := setupAuth(TEST_SECRET)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
