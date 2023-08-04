package database

import (
	"errors"
	"github.com/guneyin/bist-tools/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	db, err := gorm.Open(postgres.Open(config.Cfg.DB.URL), &gorm.Config{})
	if err != nil {
		return errors.New("failed to connect database")
	}

	DB = db

	return nil
}
