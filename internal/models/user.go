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
	OauthID       string
	Name          string
	Email         string
	AvatarURL     string
	PasswordHash  string
	EmailVerified bool
	Provider      string // "google", "local"
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

func (u *UserRepo) FindByEmailAny(email string) (*User, error) {
	user := &User{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

// UpsertLocalPassword sets a password on an existing user (account linking) or creates a new local user.
func (u *UserRepo) UpsertLocalPassword(email string, hash string, name string) (*User, error) {
	existing := &User{}
	err := db.Where("email = ?", email).First(existing).Error
	if err == nil {
		// Account exists (could be Google) — just add the password
		existing.PasswordHash = hash
		existing.EmailVerified = true
		return existing, db.Save(existing).Error
	}
	// No existing account — create a new local one
	user := &User{
		Name:          name,
		Email:         email,
		PasswordHash:  hash,
		EmailVerified: true,
		Provider:      "local",
	}
	return user, db.Create(user).Error
}

func (u *UserRepo) Delete(userID uint) error {
	db.Where("user_id = ?", userID).Delete(&Attendance{})
	db.Where("creator_user_id = ?", userID).Delete(&Event{})
	return db.Delete(&User{}, userID).Error
}

func (u *UserRepo) UpdatePassword(userID uint, hash string) error {
	return db.Model(&User{}).Where("id = ?", userID).Update("password_hash", hash).Error
}

func (u *UserRepo) FindByEmailWithPassword(email string) (*User, error) {
	user := &User{}
	err := db.Where("email = ? AND password_hash != ''", email).First(user).Error
	return user, err
}
