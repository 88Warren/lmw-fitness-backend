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

	for _, program := range programmes {
		var existingProgram models.WorkoutProgram
		if err := DB.Where("name = ?", program.Name).First(&existingProgram).Error; err != nil {
			if err := DB.Create(&program).Error; err != nil {
				log.Printf("Failed to create program %s: %v", program.Name, err)
				continue
			}
			log.Printf("Created program: %s", program.Name)
		} else {
			program = existingProgram
			log.Printf("Program already exists: %s", program.Name)
		}
	}
}
