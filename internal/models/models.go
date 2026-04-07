package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func ConnectDB(dsn string) {
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Event{})
	db.AutoMigrate(&Attendance{})
	db.AutoMigrate(&EmailVerificationToken{})
}
