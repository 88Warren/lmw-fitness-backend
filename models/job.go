// models/job.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	SessionID     string
	CustomerEmail string
	Status        string
	Attempts      int
	LastAttempt   time.Time
}
