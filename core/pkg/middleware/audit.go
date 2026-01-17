package middleware

import (
	"backend/core/pkg/audit"
	"backend/core/pkg/errorsx"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func Audit(audit *audit.Audit) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		status := c.Response().StatusCode()
		if err != nil && status == http.StatusOK {
			var fe *fiber.Error
			if errors.As(err, &fe) {
				status = fe.Code
			} else {
				status = http.StatusInternalServerError
			}
		}

		if err == nil && c.Path() == "/api" {
			return nil
		}

		if err == nil && c.Path() == "/api/v1/maintenance/health" {
			return nil
		}

		headers := make(map[string]string)
		for _, kv := range c.Request().Header.All() {
			key := string(kv[0])
			val := string(kv[1])
			headers[key] = val
		}

		meta := map[string]interface{}{
			"request_id": c.Get("X-Request-ID", ""),
			"ip":         c.IP(),
			"status":     status,
			"method":     c.Method(),
			"path":       c.Path(),
			"query":      c.Context().QueryArgs().String(),
			"user_agent": c.Get("User-Agent"),
			"latency":    time.Since(start).Milliseconds(),
			"response":   string(c.Response().Body()),
			"payload":    string(c.Body()),
			"headers":    headers,
		}

		level := "INFO"
		switch {
		case status >= http.StatusInternalServerError:
			level = "ERROR"
		case status >= http.StatusBadRequest:
			level = "WARN"
		}

		audit.Log(level, "HTTP request", meta, nil)

		return errorsx.Wrap(err, "Error handling request")
	}
}
