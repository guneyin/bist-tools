package migration

import (
	"github.com/guneyin/bist-tools/internal/repo/transaction"
	"github.com/guneyin/bist-tools/internal/repo/user"
	"github.com/guneyin/bist-tools/pkg/database"
)

func Init() error {
	db := database.DB

	return db.AutoMigrate(&transaction.Transaction{}, &user.User{})
}
