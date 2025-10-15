package tests

import (
	"os"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/database"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	os.Setenv("GO_ENV", "test")

	// Setup test environment variables
	SetupTestEnvironment()

	// Use pipeline environment variables if available, otherwise use localhost for local testing
	if os.Getenv("DB_HOST") == "" {
		// Check if we're in CI/CD pipeline or local development
		if os.Getenv("CI") != "" || os.Getenv("HARNESS_BUILD_ID") != "" {
			os.Setenv("DB_HOST", "database")
		} else {
			os.Setenv("DB_HOST", "localhost")
		}
	}
	if os.Getenv("DB_USER") == "" {
		os.Setenv("DB_USER", "postgres")
	}
	if os.Getenv("DB_PASSWORD") == "" {
		os.Setenv("DB_PASSWORD", "password")
	}
	if os.Getenv("DB_NAME") == "" {
		os.Setenv("DB_NAME", "testdb")
	}
	if os.Getenv("DB_PORT") == "" {
		os.Setenv("DB_PORT", "5432")
	}
	if os.Getenv("DB_SSLMODE") == "" {
		os.Setenv("DB_SSLMODE", "disable")
	}

	database.InitLogger()
	defer database.SyncLogger()

	config.LoadEnv()

	// Try to connect to database, but don't fail if it's not available locally
	database.ConnectToDB()
	testDB = database.GetDB()

	if testDB != nil {
		database.MigrateDB()
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}

func cleanup() {
	if testDB != nil {
		testDB.Exec("TRUNCATE TABLE jobs, users, newsletters, blogs, workout_blocks, workouts, programs CASCADE")
	}
}

func GetTestDB() *gorm.DB {
	return testDB
}
