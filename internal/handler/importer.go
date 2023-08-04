package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/internal/repo/importer"
)

func ImporterImport(c *fiber.Ctx) error {
	b := c.Query("broker")

	ts, err := importer.Import(c, b)
	if err != nil {
		return middleware.RaiseHTTPError(c, err)
	}

	return c.JSON(ts)
}

func ImporterApply(c *fiber.Ctx) error {
	err := importer.Apply(c)
	if err != nil {
		return middleware.RaiseHTTPError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
