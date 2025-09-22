package models

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutProgram struct {
	gorm.Model
	Name        string       `gorm:"not null" json:"name"`
	Description string       `json:"description"`
	Difficulty  string       `gorm:"not null" json:"difficulty"`
	Duration    int          `json:"duration"`
	IsActive    bool         `gorm:"default:true" json:"isActive"`
	Days        []WorkoutDay `gorm:"foreignKey:ProgramID" json:"days"`
}

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

type WorkoutBlock struct {
	gorm.Model
	DayID       uint              `gorm:"not null" json:"dayId"`
	BlockType   string            `gorm:"not null" json:"blockType"`
	BlockRounds int               `json:"blockRounds"`
	RoundRest   string            `json:"roundRest"`
	BlockNotes  string            `json:"blockNotes"`
	Exercises   []WorkoutExercise `gorm:"foreignKey:BlockID" json:"exercises"`
	Day         WorkoutDay        `gorm:"foreignKey:DayID" json:"-"`
}

type WorkoutExercise struct {
	gorm.Model
	BlockID       uint         `gorm:"not null" json:"blockId"`
	ExerciseID    uint         `gorm:"not null" json:"exerciseId"`
	Order         int          `gorm:"not null" json:"order"`
	Reps          string       `json:"reps"`
	Duration      string       `json:"duration"`
	WorkRestRatio string       `json:"workRestRatio"`
	Rest          string       `json:"rest"`
	Tips          string       `json:"tips"`
	Instructions  string       `json:"instructions"`
	WorkoutBlock  WorkoutBlock `gorm:"foreignKey:BlockID" json:"-"`
	Exercise      Exercise     `gorm:"foreignKey:ExerciseID" json:"exercise"`
}

type WorkoutStep struct {
	ID          uint              `json:"id"`
	StepNumber  int               `json:"stepNumber"`
	Reps        string            `json:"reps"`
	Exercises   []WorkoutExercise `json:"exercises"`
	IsCompleted bool              `json:"isCompleted"`
}

type Exercise struct {
	gorm.Model
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `json:"description"`
	Category        string    `gorm:"not null" json:"category"`
	VideoID         string    `json:"videoId"`
	Tips            string    `json:"tips"`
	Instructions    string    `json:"instructions"`
	ModificationID  *uint     `json:"modificationId"`
	Modification    *Exercise `gorm:"foreignKey:ModificationID" json:"modification"`
	ModificationID2 *uint
	Modification2   *Exercise `gorm:"foreignKey:ModificationID2"`
}

type UserWorkoutSession struct {
	gorm.Model
	UserID        uint
	WorkoutDayID  uint
	Status        string `gorm:"default:in_progress"`
	CompletedDate *time.Time
}
