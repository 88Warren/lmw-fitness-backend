package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email               string               `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	PasswordHash        string               `gorm:"not null" json:"-"`
	Role                string               `gorm:"not null;default:'user'" json:"role"`
	PasswordResetTokens []PasswordResetToken `gorm:"foreignKey:UserID"`
	MustChangePassword  bool                 `gorm:"default:true"`
	AuthTokens          []AuthToken          `gorm:"foreignKey:UserID"`
}

type UserResponse struct {
	ID                 uint   `json:"id"`
	Email              string `json:"email"`
	Role               string `json:"role"`
	MustChangePassword bool   `json:"mustChangePassword"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type PasswordResetToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID"`
}
