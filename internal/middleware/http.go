package middleware

import "github.com/gofiber/fiber/v2"

func RaiseHTTPError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error(), "data": nil})
}

func SendHTTPSuccess(c *fiber.Ctx, msg string, data any) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "success", "message": msg, "data": data})
}
