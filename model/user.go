package model

import (
	"errors"
)

type User struct {
	UserId   uint   `gorm:"primary_key"`
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Urls     []URL  `gorm:"foreignkey:user_id"`
}

func NewUser(username, password string) (*User, error) {
	if len(password) == 0 || len(username) == 0 {
		return nil, errors.New("username of password is empty")
	}
	return &User{Username: username, Password: password}, nil
}

func (user *User) ValidatePassword(pass string) bool {
	return user.Password == pass
}
