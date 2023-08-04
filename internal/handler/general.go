package handler

import "github.com/gofiber/fiber/v2"

func GeneralStatus(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "server running.."})
}
