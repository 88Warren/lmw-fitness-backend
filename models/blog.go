package models

import (
	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Date        string `json:"date" binding:"required"`
	ImageURL    string `json:"image" binding:"omitempty,url"`
	Excerpt     string `json:"excerpt" binding:"required"`
	FullContent string `json:"fullContent" binding:"required"`
}
