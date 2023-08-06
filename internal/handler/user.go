package handler

import (
	"errors"
	"github.com/google/uuid"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/internal/repo/user"
	"github.com/guneyin/bist-tools/pkg/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	return uid == n
}

func validUser(id string, p string) bool {
	u, err := user.GetByID(uuid.MustParse(id))
	if err != nil {
		return false
	}

	if u.UserName == "" {
		return false
	}

	if !CheckPasswordHash(p, u.Password) {
		return false
	}
	return true
}

var (
	errUserNotFound   = errors.New("USER_NOT_FOUND")
	errInvalidPayload = errors.New("INVALID_PAYLOAD")
)

func GetUser(c *fiber.Ctx) error {
	uid := c.Params("id")

	u, _ := user.GetByID(uuid.MustParse(uid))
	if u.UserName == "" {
		return middleware.RaiseHTTPError(c, errUserNotFound)
	}

	return middleware.SendHTTPSuccess(c, "user found", u.Safe())
}

func UserMe(c *fiber.Ctx) error {
	uid := middleware.GetUserID(c)

	u, _ := user.GetByID(uid)
	if u.UserName == "" {
		return middleware.RaiseHTTPError(c, errUserNotFound)
	}

	return middleware.SendHTTPSuccess(c, "user found", u.Safe())
}

func CreateUser(c *fiber.Ctx) error {
	u := new(user.User)
	if err := c.BodyParser(u); err != nil {
		return middleware.RaiseHTTPError(c, errInvalidPayload)
	}

	hash, err := hashPassword(u.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}
	u.Password = hash

	created, err := user.Create(u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": created.Safe()})
}

func UpdateUser(c *fiber.Ctx) error {
	type UpdateUserInput struct {
		Names string `json:"names"`
	}
	var uui UpdateUserInput
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	db := database.DB
	var user user.User

	db.First(&user, id)
	db.Save(&user)

	return c.JSON(fiber.Map{"status": "success", "message": "User successfully updated", "data": user})
}

func DeleteUser(c *fiber.Ctx) error {
	type PasswordInput struct {
		Password string `json:"password"`
	}
	var pi PasswordInput
	if err := c.BodyParser(&pi); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Review your input", "data": err})
	}
	id := c.Params("id")
	token := c.Locals("user").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})

	}

	if !validUser(id, pi.Password) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Not valid user", "data": nil})

	}

	db := database.DB
	var u user.User

	db.First(&u, id)

	db.Delete(&u)
	return c.JSON(fiber.Map{"status": "success", "message": "User successfully deleted", "data": nil})
}
