package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null;uniqueIndex"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}
