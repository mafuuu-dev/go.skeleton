package middleware

import (
	"backend/core/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors(config *config.Config) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     config.CORSOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Authorization, Content-Type, X-Requested-With",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	})
}
