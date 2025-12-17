package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func Compress(level compress.Level, excludePaths []string) fiber.Handler {
	excluded := make(map[string]struct{}, len(excludePaths))
	for _, path := range excludePaths {
		excluded[path] = struct{}{}
	}

	return compress.New(compress.Config{
		Level: level,
		Next: func(c *fiber.Ctx) bool {
			if _, ok := excluded[c.Path()]; ok {
				return true
			}

			return false
		},
	})
}
