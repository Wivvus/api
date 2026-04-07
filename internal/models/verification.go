package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	gorm.Model
	Token     string    `gorm:"uniqueIndex"`
	Email     string
	Name      string
	ExpiresAt time.Time
	Used      bool
}

type VerificationRepo struct{}

func (v *VerificationRepo) Create(email, name string) (*EmailVerificationToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	t := &EmailVerificationToken{
		Token:     hex.EncodeToString(b),
		Email:     email,
		Name:      name,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return t, db.Create(t).Error
}

func (v *VerificationRepo) FindValid(token string) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	err := db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&t).Error
	return &t, err
}

func (v *VerificationRepo) MarkUsed(id uint) error {
	return db.Model(&EmailVerificationToken{}).Where("id = ?", id).Update("used", true).Error
}
