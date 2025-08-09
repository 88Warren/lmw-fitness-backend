package models

import (
	"time"

	"gorm.io/gorm"
)

// WorkoutProgram represents a complete workout program (e.g., 30-Day Beginner)
type WorkoutProgram struct {
	gorm.Model
	Name        string       `gorm:"not null" json:"name"`
	Description string       `json:"description"`
	Difficulty  string       `gorm:"not null" json:"difficulty"` // beginner, advanced
	Duration    int          `json:"duration"`                   // number of days
	IsActive    bool         `gorm:"default:true" json:"isActive"`
	Days        []WorkoutDay `gorm:"foreignKey:ProgramID" json:"days"`
}

// WorkoutDay represents a single day's workout
type WorkoutDay struct {
	gorm.Model
	ProgramID     uint           `gorm:"not null" json:"programId"`
	DayNumber     int            `gorm:"not null" json:"dayNumber"`
	Title         string         `gorm:"not null" json:"title"`
	Description   string         `json:"description"`
	Warmup        string         `json:"warmup"`
	Cooldown      string         `json:"cooldown"`
	WorkoutBlocks []WorkoutBlock `gorm:"foreignKey:DayID" json:"workoutBlocks"`
	Program       WorkoutProgram `gorm:"foreignKey:ProgramID" json:"-"`
}

// WorkoutBlock is the new dynamic element
type WorkoutBlock struct {
	gorm.Model
	DayID       uint              `gorm:"not null" json:"dayId"`
	BlockType   string            `gorm:"not null" json:"blockType"` // e.g., "Circuit", "AMRAP", "Tabata", "EMOM", "Standard"
	BlockRounds int               `json:"blockRounds"`               // e.g., 3 rounds for a circuit
	BlockNotes  string            `json:"blockNotes"`                // e.g., "Rest 60-90 seconds between rounds"
	Exercises   []WorkoutExercise `gorm:"foreignKey:BlockID" json:"exercises"`
	Day         WorkoutDay        `gorm:"foreignKey:DayID" json:"-"`
}

// WorkoutExercise represents a single exercise within a WorkoutBlock
type WorkoutExercise struct {
	gorm.Model
	BlockID       uint         `gorm:"not null" json:"blockId"`
	ExerciseID    uint         `gorm:"not null" json:"exerciseId"`
	Order         int          `gorm:"not null" json:"order"`
	Reps          string       `json:"reps"`          // e.g., "10-12", "max reps", "15"
	Duration      string       `json:"duration"`      // e.g., "45 seconds"
	WorkRestRatio string       `json:"workRestRatio"` // e.g., "20s on, 10s off" for Tabata
	RestDuration  string       `json:"restDuration"`
	Tips          string       `json:"tips"`
	Instructions  string       `json:"instructions"`
	WorkoutBlock  WorkoutBlock `gorm:"foreignKey:BlockID" json:"-"`
	Exercise      Exercise     `gorm:"foreignKey:ExerciseID" json:"exercise"`
}

// Exercise represents a base exercise that can be used in multiple workouts
type Exercise struct {
	gorm.Model
	Name         string `gorm:"not null" json:"name"`
	Description  string `json:"description"`
	Category     string `gorm:"not null" json:"category"` // legs, upper_body, core, cardio, full_body
	VideoID      string `json:"videoId"`                  // YouTube video ID or direct URL
	Tips         string `json:"tips"`
	Instructions string `json:"instructions"`
}

type UserWorkoutSession struct {
	gorm.Model
	UserID        uint
	WorkoutDayID  uint
	Status        string `gorm:"default:in_progress"` // e.g., "in_progress", "completed"
	CompletedDate *time.Time
}
