package handler

import (
	"github.com/gofiber/fiber/v2"
	broker "github.com/guneyin/gobist-broker"
)

func BrokerList(c *fiber.Ctx) error {
	brokers := broker.GetBrokers()

	return c.JSON(brokers.ToList())
}
