package models

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	ImageURL    string `json:"image" binding:"omitempty,url"`
	Excerpt     string `json:"excerpt" binding:"required"`
	FullContent string `json:"fullContent" binding:"required"`
	IsFeatured  bool   `json:"isFeatured"`
}
