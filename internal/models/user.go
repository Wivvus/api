package models

import (
	"github.com/markbates/goth"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	NickName  string
	AvatarURL string
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

func (u *UserRepo) CreateOrUpdate(user *User) error {
	resp := db.Model(&User{}).Where("id = ?", user.ID).Updates(user)
	if resp.Error != nil {
		return resp.Error
	}
	if resp.RowsAffected == 0 {
		return db.Create(user).Error
	}
	return nil
}

func UserFromGoogleOAuth(user goth.User) *User {
	return &User{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		NickName:  user.NickName,
		AvatarURL: user.AvatarURL,
	}
}
