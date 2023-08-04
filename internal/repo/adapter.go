package repo

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type AdapterWithCache interface {
	GetAddress() string
	SaveToCache(ctx context.Context) error
	GetFromCache(ctx context.Context) error
	DeleteFromCache(ctx context.Context) error
	ToJSON() []byte
	FromJSON(d []byte) error
}

type AdapterWithDB interface {
	SaveToDB(c *fiber.Ctx) error
}

type AdapterFull interface {
	AdapterWithDB
	AdapterWithCache
}
