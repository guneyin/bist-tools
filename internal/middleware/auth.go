package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guneyin/bist-tools/pkg/config"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup:  "cookie:auth-token",
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Cfg.JWT.Secret)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func GetUserID(c *fiber.Ctx) uuid.UUID {
	u := c.Locals("user").(*jwt.Token)
	claims := u.Claims.(jwt.MapClaims)
	uid := claims["user_id"].(string)

	return uuid.MustParse(uid)
}
