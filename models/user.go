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
	MustChangePassword  bool                 `gorm:"default:false"`
	AuthTokens          []AuthToken          `gorm:"foreignKey:UserID"`
	UserPrograms        []UserProgram        `gorm:"foreignKey:UserID"`
	CompletedDays       map[string]int       `json:"completedDays" gorm:"serializer:json"`
	ProgramStartDates   map[string]time.Time `json:"programStartDates" gorm:"serializer:json"`
	CompletedDaysList   map[string][]int     `json:"completedDaysList" gorm:"serializer:json"`
	Timezone            string               `gorm:"default:'UTC'" json:"timezone"`
}

type UserResponse struct {
	ID                 uint                 `json:"id"`
	Email              string               `json:"email"`
	Role               string               `json:"role"`
	MustChangePassword bool                 `json:"mustChangePassword"`
	PurchasedPrograms  []string             `json:"purchasedPrograms"`
	CompletedDays      map[string]int       `json:"completedDays"`
	ProgramStartDates  map[string]time.Time `json:"programStartDates"`
	CompletedDaysList  map[string][]int     `json:"completedDaysList"`
	UnlockedDays       map[string]int       `json:"unlockedDays"`
	Timezone           string               `json:"timezone"`
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

type UserProgram struct {
	gorm.Model
	UserID         uint `gorm:"not null"`
	ProgramID      uint `gorm:"not null"`
	User           User
	WorkoutProgram WorkoutProgram `gorm:"foreignKey:ProgramID"`
}
