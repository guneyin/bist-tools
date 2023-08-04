package user

import (
	"github.com/guneyin/bist-tools/pkg/database"
)

type User struct {
	database.Model
	UserName string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"password,omitempty"`
}

func (u User) Safe() User {
	r := u
	r.Password = ""

	return r
}
