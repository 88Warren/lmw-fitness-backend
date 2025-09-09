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
	// Set test environment
	os.Setenv("GO_ENV", "test")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "lmwfitness_test")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_SSLMODE", "disable")

	// Initialize logger
	database.InitLogger()
	defer database.SyncLogger()

	// Load config
	config.LoadEnv()

	// Connect to test database
	database.ConnectToDB()
	testDB = database.GetDB()

	// Run migrations
	database.MigrateDB()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func cleanup() {
	if testDB != nil {
		// Clean up test data
		testDB.Exec("TRUNCATE TABLE jobs, users, newsletters, blogs, workout_blocks, workouts, programs CASCADE")
	}
}

func GetTestDB() *gorm.DB {
	return testDB
}
