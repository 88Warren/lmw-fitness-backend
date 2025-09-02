// controllers/workout_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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

type ProgramDetailsResponse struct {
	Title     string `json:"title"`
	TotalDays int    `json:"totalDays"`
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
		Preload("WorkoutBlocks.Exercises.Exercise.Modification").
		Preload("WorkoutBlocks.Exercises.Exercise").
		Preload("WorkoutBlocks.Exercises").
		First(&workoutDay).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout day not found"})
		return
	}

	c.JSON(http.StatusOK, workoutDay)
}

// GetWarmup retrieves the warmup routine for a specific program
func (wc *WorkoutController) GetWarmup(c *gin.Context) {
	wc.getRoutineByProgramName(c, "warmup")
}

// GetCooldown retrieves the cooldown routine for a specific program
func (wc *WorkoutController) GetCooldown(c *gin.Context) {
	wc.getRoutineByProgramName(c, "cooldown")
}

func (wc *WorkoutController) getRoutineByProgramName(c *gin.Context, routineType string) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	programName := c.Param("programName")

	var user models.User
	if err := wc.DB.Preload("UserPrograms.WorkoutProgram").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	isAuthorized := false
	if user.Role == "admin" {
		isAuthorized = true
	} else {
		for _, userProgram := range user.UserPrograms {
			if userProgram.WorkoutProgram.Name == programName {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorised to view this program."})
		return
	}

	var exerciseName string
	switch routineType {
	case "warmup":
		exerciseName = "Warm Up"
	case "cooldown":
		exerciseName = "Cool Down"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid routine type specified."})
		return
	}

	var exercise models.Exercise
	if err := wc.DB.Where("name = ?", exerciseName).First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": routineType + " routine not found"})
		return
	}

	// Construct the response using the exercise data from the database
	response := gin.H{
		"videoUrl":     fmt.Sprintf("https://www.youtube.com/embed/%s", exercise.VideoID),
		"description":  exercise.Description,
		"instructions": exercise.Instructions,
		"tips":         exercise.Tips,
		"category":     exercise.Category,
	}

	c.JSON(http.StatusOK, response)
}

func (wc *WorkoutController) GetWorkoutDayByProgramAndDay(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		// fmt.Println("Backend: Authorization check failed - No userID in context.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// fmt.Printf("Backend: User ID from context: %v\n", userID)

	programName := c.Param("programName")
	dayNumberStr := c.Param("dayNumber")
	// fmt.Printf("Backend: URL Params - ProgramName: %s, DayNumber: %s\n", programName, dayNumberStr)

	dayNumber, err := strconv.Atoi(dayNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day number"})
		return
	}

	// 1. Fetch the user and PRELOAD the nested WorkoutProgram data
	var user models.User
	if err := wc.DB.Preload("UserPrograms.WorkoutProgram").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}
	// fmt.Printf("Backend: User from DB: Email: %s, Role: %s\n", user.Email, user.Role)
	// fmt.Printf("Backend: User purchased programs: %v\n", user.UserPrograms)

	// 2. Perform the authorization check
	isAuthorized := false
	if user.Role == "admin" {
		isAuthorized = true
	} else {
		// Iterate through the PRELOADED UserPrograms and check the WorkoutProgram's Name
		for _, userProgram := range user.UserPrograms {
			if userProgram.WorkoutProgram.Name == programName {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorised to view this program."})
		return
	}

	// 3. If authorized, proceed to fetch and return the workout data
	var program models.WorkoutProgram
	if err := wc.DB.Where("name = ?", programName).First(&program).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
		return
	}

	var workoutDay models.WorkoutDay
	if err := wc.DB.Where("program_id = ? AND day_number = ?", program.ID, dayNumber).
		Preload("WorkoutBlocks").
		Preload("WorkoutBlocks.Exercises").
		Preload("WorkoutBlocks.Exercises.Exercise").
		Preload("WorkoutBlocks.Exercises.Exercise.Modification").
		First(&workoutDay).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout day not found for this program"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Exercise completion recorded"})
}

// CompleteWorkoutDay marks an entire workout session as complete
func (wc *WorkoutController) CompleteWorkoutDay(c *gin.Context) {
	// Get user from middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ProgramName string `json:"programName" binding:"required"`
		DayNumber   int    `json:"dayNumber" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user
	var user models.User
	if err := wc.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update user's completed days for this program
	// Initialize completedDays if nil
	if user.CompletedDays == nil {
		user.CompletedDays = make(map[string]int)
	}

	// Update the completed day for this program (only if it's higher than current)
	currentCompleted := user.CompletedDays[req.ProgramName]
	if req.DayNumber > currentCompleted {
		user.CompletedDays[req.ProgramName] = req.DayNumber
	}

	// Save the user with updated completion data
	if err := wc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user completion status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Workout day completed successfully",
		"completedDay": req.DayNumber,
		"programName":  req.ProgramName,
	})
}

// GetUserProgress retrieves a summary of the user's completed workouts
func (wc *WorkoutController) GetUserProgress(c *gin.Context) {
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

func (wc *WorkoutController) GetProgramList(c *gin.Context) {
	programName := c.Param("programID")

	var program models.WorkoutProgram
	if err := wc.DB.Where("name = ?", programName).First(&program).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve program details"})
		return
	}

	var totalDays int64
	if err := wc.DB.Model(&models.WorkoutDay{}).Where("program_id = ?", program.ID).Count(&totalDays).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count total days for program"})
		return
	}

	response := ProgramDetailsResponse{
		Title:     program.Name,
		TotalDays: int(totalDays),
	}

	c.JSON(http.StatusOK, response)
}
