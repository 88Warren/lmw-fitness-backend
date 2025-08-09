package database

import (
	"log"

	"github.com/88warren/lmw-fitness-backend/models"
)

func ProgramSeed() {

	log.Println("Seeding programme data...")

	programmes := []models.WorkoutProgram{
		{
			Name:        "30-Day Beginner Program",
			Description: "A comprehensive 30-day program designed for fitness beginners",
			Difficulty:  "beginner",
			Duration:    30,
			IsActive:    true,
		},
		{
			Name:        "30-Day Advanced Program",
			Description: "A challenging 30-day program for experienced fitness enthusiasts",
			Difficulty:  "advanced",
			Duration:    30,
			IsActive:    true,
		},
	}

	for _, p := range programmes {
		var existingProgram models.WorkoutProgram
		if err := DB.Where("name = ?", p.Name).First(&existingProgram).Error; err == nil {
			log.Printf("Program '%s' (%s) already exists", existingProgram.Name, existingProgram.Difficulty)
			continue
		}
		if err := DB.Create(&p).Error; err != nil {
			log.Printf("Failed to create program '%s' (%s): %v", p.Name, p.Difficulty, err)
		} else {
			log.Printf("Successfully created program '%s' (%s).", p.Name, p.Difficulty)
		}
	}
}
