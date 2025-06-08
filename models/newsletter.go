package models

import (
	"time"

	"gorm.io/gorm"
)

type NewsletterSubscriber struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null"`
	IsConfirmed  bool   `gorm:"default:false"`
	ConfirmToken string
	SubscribedAt time.Time
}
