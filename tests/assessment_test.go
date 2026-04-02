package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetDay1Assessment(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Setup router
	router := config.SetupServer()
	assessmentController := controllers.NewAssessmentController(db)
	routes.RegisterAssessmentRoutes(router, assessmentController)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}
	db.Create(&user)

	// Create test exercise
	exercise := models.Exercise{
		Name:        "Press Ups",
		Description: "Test exercise",
	}
	db.Create(&exercise)

	// Create Day 1 assessment
	reps := 25
	assessment := models.FitnessAssessment{
		UserID:       user.ID,
		ProgramName:  "beginner-program",
		DayNumber:    1,
		ExerciseID:   exercise.ID,
		ExerciseName: exercise.Name,
		Reps:         &reps,
	}
	db.Create(&assessment)

	// Get JWT token for authentication
	token := createTestJWTToken(user.ID)

	// Test getting Day 1 assessment
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/assessments/beginner-program/exercise/%d/day1", exercise.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response controllers.AssessmentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, assessment.ProgramName, response.ProgramName)
	assert.Equal(t, assessment.DayNumber, response.DayNumber)
	assert.Equal(t, assessment.ExerciseID, response.ExerciseID)
	assert.Equal(t, *assessment.Reps, *response.Reps)

	// Cleanup
	db.Delete(&assessment)
	db.Delete(&exercise)
	db.Delete(&user)
}

func TestGetDay1AssessmentNotFound(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Setup router
	router := config.SetupServer()
	assessmentController := controllers.NewAssessmentController(db)
	routes.RegisterAssessmentRoutes(router, assessmentController)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}
	db.Create(&user)

	// Create test exercise (but no Day 1 assessment)
	exercise := models.Exercise{
		Name:        "Press Ups",
		Description: "Test exercise",
	}
	db.Create(&exercise)

	// Get JWT token for authentication
	token := createTestJWTToken(user.ID)

	// Test getting non-existent Day 1 assessment
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/assessments/beginner-program/exercise/%d/day1", exercise.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Day 1 assessment not found")

	// Cleanup
	db.Delete(&exercise)
	db.Delete(&user)
}

func TestSaveAssessmentDay30(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Setup router
	router := config.SetupServer()
	assessmentController := controllers.NewAssessmentController(db)
	routes.RegisterAssessmentRoutes(router, assessmentController)

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}
	db.Create(&user)

	// Create test exercise
	exercise := models.Exercise{
		Name:        "Press Ups",
		Description: "Test exercise",
	}
	db.Create(&exercise)

	// Get JWT token for authentication
	token := createTestJWTToken(user.ID)

	// Test saving Day 30 assessment
	reps := 35
	assessmentData := controllers.SaveAssessmentRequest{
		ProgramName:  "beginner-program",
		DayNumber:    30,
		ExerciseID:   exercise.ID,
		ExerciseName: exercise.Name,
		Reps:         &reps,
		Notes:        "Improved from Day 1!",
	}

	jsonData, _ := json.Marshal(assessmentData)
	req, _ := http.NewRequest("POST", "/api/assessments/save", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Assessment saved successfully", response["message"])

	// Cleanup
	db.Where("user_id = ? AND exercise_id = ?", user.ID, exercise.ID).Delete(&models.FitnessAssessment{})
	db.Delete(&exercise)
	db.Delete(&user)
}

// Helper function to create JWT token for testing
func createTestJWTToken(userID uint) string {
	// This is a simplified token for testing
	// In a real implementation, you'd use the same JWT creation logic as your auth system
	return "test-jwt-token"
}
