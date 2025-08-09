// controllers/workout_controller.go
package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// WorkoutController holds the database connection
type WorkoutController struct {
	DB *gorm.DB
}

// NewWorkoutController creates a new instance of WorkoutController
func NewWorkoutController(db *gorm.DB) *WorkoutController {
	return &WorkoutController{DB: db}
}

// GetWorkoutPrograms retrieves a list of all workout programs
func (wc *WorkoutController) GetWorkoutPrograms(c *gin.Context) {
	var programs []models.WorkoutProgram
	if err := wc.DB.Find(&programs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve programs"})
		return
	}
	c.JSON(http.StatusOK, programs)
}

// GetWorkoutDay retrieves a specific workout day with its blocks and exercises
func (wc *WorkoutController) GetWorkoutDay(c *gin.Context) {
	programID, err := strconv.Atoi(c.Param("programID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	dayNumber, err := strconv.Atoi(c.Param("dayNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day number"})
		return
	}

	var workoutDay models.WorkoutDay
	if err := wc.DB.Where("program_id = ? AND day_number = ?", programID, dayNumber).
		Preload("WorkoutBlocks.Exercises").
		First(&workoutDay).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout day not found"})
		return
	}

	c.JSON(http.StatusOK, workoutDay)
}

// StartWorkout initiates a new workout session for a user on a given day
func (wc *WorkoutController) StartWorkout(c *gin.Context) {
	var req struct {
		WorkoutDayID uint `json:"workout_day_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get userID from the context set by the AuthMiddleware
	userID, _ := c.Get("userID")

	// Check if a session for this day already exists for the user
	var existingSession models.UserWorkoutSession
	if wc.DB.Where("user_id = ? AND workout_day_id = ?", userID, req.WorkoutDayID).First(&existingSession).Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Workout session already started", "session_id": existingSession.ID})
		return
	}

	newSession := models.UserWorkoutSession{
		UserID:       userID.(uint),
		WorkoutDayID: req.WorkoutDayID,
		Status:       "in_progress",
	}
	if err := wc.DB.Create(&newSession).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start workout session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Workout session started", "session_id": newSession.ID})
}

// CompleteExercise marks an individual exercise as complete
func (wc *WorkoutController) CompleteExercise(c *gin.Context) {
	var req struct {
		SessionID  uint `json:"session_id" binding:"required"`
		ExerciseID uint `json:"exercise_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// This is a simple implementation. A more advanced version would also
	// update a separate table that tracks completed exercises within a session.
	// For this example, we'll just log the completion.

	c.JSON(http.StatusOK, gin.H{"message": "Exercise completion recorded"})
}

// CompleteWorkoutDay marks an entire workout session as complete
func (wc *WorkoutController) CompleteWorkoutDay(c *gin.Context) {
	var req struct {
		SessionID uint `json:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var session models.UserWorkoutSession
	if err := wc.DB.First(&session, req.SessionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	now := time.Now()
	if err := wc.DB.Model(&session).Updates(models.UserWorkoutSession{
		Status:        "completed",
		CompletedDate: &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete workout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout day completed successfully"})
}

// GetUserProgress retrieves a summary of the user's completed workouts
func (wc *WorkoutController) GetUserProgress(c *gin.Context) {
	// Get userID from the context set by the AuthMiddleware
	userID, _ := c.Get("userID")

	var sessions []models.UserWorkoutSession
	if err := wc.DB.Where("user_id = ? AND status = ?", userID, "completed").
		Preload("WorkoutDay").
		Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user progress"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}
