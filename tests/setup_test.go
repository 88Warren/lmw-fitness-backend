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
		// Clean up tables that exist, ignore errors for missing tables
		tables := []string{"jobs", "users", "newsletters", "blogs", "workout_blocks", "workouts", "programs"}

		for _, table := range tables {
			// Check if table exists before trying to truncate
			var exists bool
			err := testDB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", table).Scan(&exists).Error
			if err == nil && exists {
				testDB.Exec("TRUNCATE TABLE " + table + " CASCADE")
			}
		}
	}
}

func GetTestDB() *gorm.DB {
	return testDB
}
