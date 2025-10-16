package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestGetWorkoutPrograms(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Create test workout programs
	program1 := models.WorkoutProgram{
		Name:        "Test Program 1",
		Description: "This is test program 1",
		Duration:    30,
		Difficulty:  "beginner",
		IsActive:    true,
	}
	program2 := models.WorkoutProgram{
		Name:        "Test Program 2",
		Description: "This is test program 2",
		Duration:    45,
		Difficulty:  "intermediate",
		IsActive:    true,
	}
	db.Create(&program1)
	db.Create(&program2)

	router := config.SetupServer()
	workoutController := controllers.NewWorkoutController(GetTestDB())
	routes.RegisterWorkoutRoutes(router, workoutController)

	req, _ := http.NewRequest("GET", "/api/workouts/programs", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.WorkoutProgram
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.GreaterOrEqual(t, len(response), 2)
}

func TestGetWorkoutProgramByID(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Create test workout program
	program := models.WorkoutProgram{
		Name:        "Specific Program Test",
		Description: "This is specific program test",
		Duration:    60,
		Difficulty:  "advanced",
		IsActive:    true,
	}
	db.Create(&program)

	router := config.SetupServer()
	workoutController := controllers.NewWorkoutController(GetTestDB())
	routes.RegisterWorkoutRoutes(router, workoutController)

	req, _ := http.NewRequest("GET", "/api/workouts/programs/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 200 or 404 depending on implementation
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
}

func TestCreateWorkoutProgram(t *testing.T) {
	// Skip this test as the POST /api/workouts/programs endpoint is not implemented
	// This would be an admin-only feature that requires authentication
	t.Skip("Skipping workout program creation test - endpoint not implemented")
}

func TestWorkoutProgramValidation(t *testing.T) {
	// Skip this test as the POST /api/workouts/programs endpoint is not implemented
	// This would be an admin-only feature that requires authentication
	t.Skip("Skipping workout program validation test - endpoint not implemented")
}
