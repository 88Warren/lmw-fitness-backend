package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Role         string `gorm:"not null;default:'user'" json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
