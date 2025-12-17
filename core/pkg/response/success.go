package response

import (
	"github.com/gofiber/fiber/v2"
)

type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(SuccessResponse{
		Data:    data,
		Success: true,
	})
}
