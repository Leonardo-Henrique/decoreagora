package utils

import "github.com/gofiber/fiber/v2"

func GetCurrentUserID(c *fiber.Ctx) int {
	if userID, ok := c.Locals("userID").(int); ok {
		return userID
	}
	return 0
}
