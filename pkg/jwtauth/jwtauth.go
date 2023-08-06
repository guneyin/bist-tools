package jwtauth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/guneyin/bist-tools/pkg/config"
	"time"
)

func CreateToken(id, name string) (string, error) {
	gen := jwt.New(jwt.SigningMethodHS256)

	claims := gen.Claims.(jwt.MapClaims)
	claims["username"] = name
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(config.Cfg.Duration)).Unix()

	return gen.SignedString([]byte(config.Cfg.Secret))
}
