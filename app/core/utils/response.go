package utils

import (
	"github.com/gofiber/fiber/v2"
)

func ErrorResponse(c *fiber.Ctx, statuscode int, err error) error {
	return c.Status(statuscode).JSON(fiber.Map{
		"error": err.Error(),
	})
}
