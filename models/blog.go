package models

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title       string    `json:"title"`
	Slug        string    `json:"slug" gorm:"uniqueIndex"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	ImageURL    string    `json:"image_url"`
	PublishedAt time.Time `json:"published_at"`
	IsPublished bool      `json:"is_published"`
	// UserID      uint      `json:"user_id" form:"user_id"`
	// User        User      `json:"user" form:"user" gorm:"foreignKey:UserID"`
}
