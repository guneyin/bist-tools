package importer

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/bist-tools/internal/repo/transaction"
	broker "github.com/guneyin/gobist-broker"
	"github.com/guneyin/gobist-broker/entity"
	"io"
)

func Import(c *fiber.Ctx, broker string) (*transaction.Transactions, error) {
	switch entity.EnumBroker(broker) {
	case entity.Garanti:
		return garantiImport(c)
	case entity.NCM:
		return ncmImport(c)
	default:
		return nil, errors.New("invalid broker")
	}
}

func Apply(c *fiber.Ctx) error {
	imp, err := NewImportSession(c)
	if err != nil {
		return err
	}

	return imp.SaveToDB(c)
}

func import_(c *fiber.Ctx, eb entity.EnumBroker, content []byte) (*transaction.Transactions, error) {
	b := broker.GetBroker(eb)

	ts, err := b.Parse(content)
	if err != nil {
		return nil, err
	}

	imp, err := NewImportSession(c)
	if err != nil {
		return nil, err
	}

	return imp.process(c, b, ts)
}

func garantiImport(c *fiber.Ctx) (*transaction.Transactions, error) {
	fh, err := c.FormFile("document")
	if err != nil {
		return nil, err
	}

	file, err := fh.Open()
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return import_(c, entity.Garanti, content)
}

func ncmImport(c *fiber.Ctx) (*transaction.Transactions, error) {
	content := c.Body()

	return import_(c, entity.NCM, content)
}
