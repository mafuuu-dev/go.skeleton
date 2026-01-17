package middleware

import (
	"backend/core/constants"
	pkgjwt "backend/core/pkg/jwt"
	"backend/core/pkg/scope"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

type AuthProvider interface {
	IsValidToken(c *fiber.Ctx) bool
}

func Auth(scope *scope.Scope, provider AuthProvider) fiber.Handler {
	jwtWrapper := func(c *fiber.Ctx, jwtMiddleware fiber.Handler) error {
		authorization := c.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":  string(constants.UnauthorizedError),
				"code":   http.StatusUnauthorized,
				"status": false,
			})
		}

		token := strings.TrimPrefix(authorization, "Bearer ")
		claims, err := pkgjwt.NewSecurity(scope.Config).Decode(token)
		if err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
			signed, _ := token.SignedString([]byte(scope.Config.JWTSecret))

			c.Request().Header.Set("Authorization", "Bearer "+signed)
		}

		return jwtMiddleware(c)
	}

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: []byte(scope.Config.JWTSecret),
		ContextKey: "user",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error":  string(constants.UnauthorizedError),
				"code":   http.StatusUnauthorized,
				"status": false,
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			c.Locals("player_id", claims["user_id"])
			c.Locals("account_id", claims["account_id"])
			c.Locals("currency_code", claims["currency_code"])

			if !provider.IsValidToken(c) {
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"error":  string(constants.UnauthorizedError),
					"code":   http.StatusUnauthorized,
					"status": false,
				})
			}

			return c.Next()
		},
	})

	return func(c *fiber.Ctx) error {
		return jwtWrapper(c, jwtMiddleware)
	}
}
