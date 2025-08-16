package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// log.Printf("Attempting to connect to database at %s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dns), &gorm.Config{})
		if err == nil {
			// log.Println("Database connection established")
			return
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Second * 5)
		}
	}

	log.Fatalf("Failed to connect to database after %d attempts: %v", maxRetries, err)
}

func GetDB() *gorm.DB {
	return DB
}

// MigrateDB creates all database tables
func MigrateDB() {
	log.Println("Starting database migration...")

	err := DB.AutoMigrate(
		&models.AuthToken{},
		&models.Job{},
		&models.User{},
		&models.UserProgram{},
		&models.PasswordResetToken{},
		&models.WorkoutProgram{},
		&models.WorkoutDay{},
		&models.WorkoutBlock{},
		&models.Exercise{},
		&models.WorkoutExercise{},
		&models.UserWorkoutSession{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")
}
