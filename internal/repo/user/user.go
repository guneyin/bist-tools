package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/guneyin/bist-tools/pkg/database"
	"gorm.io/gorm"
)

func GetByID(uid uuid.UUID) (*User, error) {
	db := database.DB
	var user User
	if err := db.Where(&User{}).Find(&user).Where("id = ?", uid.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetByEmail(e string) (*User, error) {
	db := database.DB
	var user User
	if err := db.Where(&User{Email: e}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func GetByUserName(u string) (*User, error) {
	db := database.DB
	var user User
	if err := db.Where(&User{UserName: u}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func Create(u *User) (*User, error) {
	db := database.DB

	u.ID = uuid.New()
	err := db.Create(u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}
