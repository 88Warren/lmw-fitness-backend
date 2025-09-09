package database

import (
	"log"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserSeed(db *gorm.DB) {
	seedAdminUser(db)
	seedGenericUser(db)
}

func seedAdminUser(db *gorm.DB) {
	adminEmail := "admin@lmwfitness.co.uk"
	adminPassword := "AdminPassword123!"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	adminUser := models.User{
		Email:              adminEmail,
		PasswordHash:       string(hashedPassword),
		Role:               "admin",
		MustChangePassword: false,
		CompletedDays:      make(map[string]int),
		ProgramStartDates:  make(map[string]time.Time),
		CompletedDaysList:  make(map[string][]int),
	}

	var existingUser models.User
	result := db.Where("email = ?", adminUser.Email).First(&existingUser)

	switch result.Error {
	case gorm.ErrRecordNotFound:
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("Failed to create admin user: %v", err)
			return
		}
		db.Model(&adminUser).Update("must_change_password", false)
		// log.Printf("Admin user '%s' created successfully.", adminUser.Email)
	case nil:
		existingUser.PasswordHash = adminUser.PasswordHash
		existingUser.Role = "admin"
		existingUser.MustChangePassword = false
		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update admin user '%s': %v", existingUser.Email, err)
			return
		}
		db.Model(&existingUser).Update("must_change_password", false)
		// log.Printf("Admin user '%s' updated successfully.", existingUser.Email)
	default:
		log.Printf("Database error checking for admin user: %v", result.Error)
		return
	}
}

func seedGenericUser(db *gorm.DB) {
	userEmail := "user@example.com"
	userPassword := "UserPassword123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash user password: %v", err)
		return
	}

	genericUser := models.User{
		Email:              userEmail,
		PasswordHash:       string(hashedPassword),
		Role:               "user",
		MustChangePassword: false,
		CompletedDays: map[string]int{
			"beginner": 3,
			"advanced": 1,
		},
		ProgramStartDates: map[string]time.Time{
			"beginner": time.Now().AddDate(0, 0, -3), // Started 3 days ago
			"advanced": time.Now().AddDate(0, 0, -1), // Started 1 day ago
		},
		CompletedDaysList: map[string][]int{
			"beginner": {1, 2, 3},
			"advanced": {1},
		},
	}

	var existingUser models.User
	result := db.Where("email = ?", genericUser.Email).First(&existingUser)

	var userID uint
	switch result.Error {
	case gorm.ErrRecordNotFound:
		if err := db.Create(&genericUser).Error; err != nil {
			log.Printf("Failed to create generic user: %v", err)
			return
		}
		db.Model(&genericUser).Update("must_change_password", false)
		userID = genericUser.ID
		// log.Printf("Generic user '%s' created successfully.", genericUser.Email)
	case nil:
		existingUser.PasswordHash = genericUser.PasswordHash
		existingUser.Role = "user"
		existingUser.MustChangePassword = false
		existingUser.CompletedDays = genericUser.CompletedDays
		existingUser.ProgramStartDates = genericUser.ProgramStartDates
		existingUser.CompletedDaysList = genericUser.CompletedDaysList
		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update generic user '%s': %v", existingUser.Email, err)
			return
		}
		db.Model(&existingUser).Update("must_change_password", false)
		userID = existingUser.ID
		// log.Printf("Generic user '%s' updated successfully.", existingUser.Email)
	default:
		log.Printf("Database error checking for generic user: %v", result.Error)
		return
	}
	seedUserPrograms(db, userID)
}

func seedUserPrograms(db *gorm.DB, userID uint) {
	var beginnerProgram models.WorkoutProgram
	var advancedProgram models.WorkoutProgram

	if err := db.Where("name = ?", "beginner-program").First(&beginnerProgram).Error; err != nil {
		log.Printf("Failed to find beginner program: %v", err)
		return
	}

	if err := db.Where("name = ?", "advanced-program").First(&advancedProgram).Error; err != nil {
		log.Printf("Failed to find advanced program: %v", err)
		return
	}

	programs := []models.UserProgram{
		{
			UserID:    userID,
			ProgramID: beginnerProgram.ID,
		},
		{
			UserID:    userID,
			ProgramID: advancedProgram.ID,
		},
	}

	for _, userProgram := range programs {
		var existingUserProgram models.UserProgram
		result := db.Where("user_id = ? AND program_id = ?", userProgram.UserID, userProgram.ProgramID).First(&existingUserProgram)

		switch result.Error {
		case gorm.ErrRecordNotFound:
			if err := db.Create(&userProgram).Error; err != nil {
				log.Printf("Failed to create user program association: %v", err)
			} else {
				log.Printf("Successfully associated user with program ID %d", userProgram.ProgramID)
			}
		case nil:
			log.Printf("User program association already exists for program ID %d", userProgram.ProgramID)
		default:
			log.Printf("Database error checking user program association: %v", result.Error)
		}
	}
}
