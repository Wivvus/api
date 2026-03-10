package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	UserExists = errors.New("User already exists")
)

type User struct {
	gorm.Model
	OauthID   string
	Name      string
	Email     string
	AvatarURL string
}

type UserAPIDecorator struct {
	Name      string
	Email     string
	AvatarURL string
	ID        string
}

func (u *User) ToAPI() *UserAPIDecorator {
	return &UserAPIDecorator{
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		ID:        u.OauthID,
	}
}

type UserRepo struct {
}

func (u *UserRepo) FindByID(id string) *User {
	user := &User{}
	db.First(user, id)
	return user
}

func (u *UserRepo) FindByEmail(email string) *User {
	user := &User{}
	db.First(user, "email = ?", email)
	return user
}

func (u *UserRepo) Create(user *User) error {
	resp := db.Model(&User{}).Where("oauth_id = ?", user.OauthID)
	if resp.Error != nil {
		return resp.Error
	}
	if resp.RowsAffected > 0 {
		return UserExists
	}
	if resp.RowsAffected == 0 {
		return db.Create(user).Error
	}
	return nil
}
