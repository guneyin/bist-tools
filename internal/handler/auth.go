package handler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/internal/repo/user"
	"github.com/guneyin/bist-tools/pkg/config"
	"github.com/guneyin/bist-tools/pkg/jwtauth"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"time"
)

var (
	errErrorOnLoginRequest = errors.New("ERROR_ON_LOGIN_REQUEST")
	errInvalidPassword     = errors.New("INVALID_PASSWORD")
)

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	li := new(LoginInput)
	if err := c.BodyParser(&li); err != nil {
		return middleware.RaiseHTTPError(c, errErrorOnLoginRequest)
	}

	identity := li.Identity
	pass := li.Password
	usr, err := new(user.User), *new(error)

	if isEmail(identity) {
		usr, err = user.GetByEmail(identity)
	} else {
		usr, err = user.GetByUserName(identity)
	}

	if usr == nil {
		return middleware.RaiseHTTPError(c, err)
	}

	if !CheckPasswordHash(pass, usr.Password) {
		return middleware.RaiseHTTPError(c, errInvalidPassword)
	}

	token, err := jwtauth.CreateToken(usr.ID.String(), usr.UserName)
	if err != nil {
		return middleware.RaiseHTTPError(c, err)
	}

	c.Cookie(&fiber.Cookie{
		Name:        "auth-token",
		Value:       token,
		Path:        "",
		Domain:      "",
		MaxAge:      0,
		Expires:     time.Now().Add(time.Hour * time.Duration(config.Cfg.Duration)),
		Secure:      true,
		HTTPOnly:    true,
		SameSite:    "",
		SessionOnly: false,
	})

	return middleware.SendHTTPSuccess(c, "login successful", nil)
}
