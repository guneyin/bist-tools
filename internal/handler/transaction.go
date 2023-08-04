package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/internal/repo/importer"
)

func TransactionsGet(c *fiber.Ctx) error {
	s, err := importer.NewImportSession(c)
	if err != nil {
		return middleware.RaiseHTTPError(c, err)
	}

	return c.JSON(s)
}
