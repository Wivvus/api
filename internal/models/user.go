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
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (u *User) ToAPI() *UserAPIDecorator {
	return &UserAPIDecorator{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
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
	existing := &User{}
	resp := db.Where("oauth_id = ?", user.OauthID).First(existing)
	if resp.Error == nil {
		user.Model = existing.Model
		return UserExists
	}
	return db.Create(user).Error
}
