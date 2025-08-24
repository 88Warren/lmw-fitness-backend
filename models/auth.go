package models

import "gorm.io/gorm"

type AuthToken struct {
	gorm.Model
	UserID      uint   `gorm:"index"`
	Token       string `gorm:"unique;size:64"`
	ProgramName string `json:"program_name"`
	DayNumber   int    `json:"day_number"`
	IsUsed      bool   `gorm:"default:false"`
	User        User   `gorm:"foreignKey:UserID"`
	SessionID   string `json:"session_id"`
}
