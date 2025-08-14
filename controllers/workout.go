// controllers/workout_controller.go
package controllers

import (
	"fmt"
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

// GetWorkoutDayByProgramAndDay retrieves a specific workout day with its blocks and exercises
func (wc *WorkoutController) GetWorkoutDayByProgramAndDay(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		fmt.Println("Backend: Authorization check failed - No userID in context.")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	fmt.Printf("Backend: User ID from context: %v\n", userID)

	programName := c.Param("programName")
	dayNumberStr := c.Param("dayNumber")
	fmt.Printf("Backend: URL Params - ProgramName: %s, DayNumber: %s\n", programName, dayNumberStr)

	dayNumber, err := strconv.Atoi(dayNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day number"})
		return
	}

	// 1. Fetch the user with their AuthTokens preloaded
	var user models.User
	fmt.Printf("Backend: Attempting to fetch user %v with AuthTokens from DB.\n", userID)
	if err := wc.DB.Preload("AuthTokens").First(&user, userID).Error; err != nil {
		fmt.Printf("Backend: Error fetching user: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	fmt.Printf("Backend: User from DB: Email: %s, Role: %s\n", user.Email, user.Role)
	fmt.Printf("Backend: Number of AuthTokens found for user: %d\n", len(user.AuthTokens))

	// Enhanced debugging of auth tokens
	for i, token := range user.AuthTokens {
		fmt.Printf("Backend: Token %d details:\n", i)
		fmt.Printf("  - ID: %d\n", token.ID)
		fmt.Printf("  - UserID: %d\n", token.UserID)
		fmt.Printf("  - ProgramName: '%s'\n", token.ProgramName)
		fmt.Printf("  - DayNumber: %d\n", token.DayNumber)
		fmt.Printf("  - IsUsed: %v\n", token.IsUsed)
		fmt.Printf("  - SessionID: %s\n", token.SessionID)
	}

	// 2. Perform the authorization check
	isAuthorized := false
	authReason := ""

	if user.Role == "admin" {
		isAuthorized = true
		authReason = "User is admin"
		fmt.Println("Backend: User is an admin, granting access.")
	} else {
		fmt.Println("Backend: User is not an admin. Checking purchased programs...")
		for _, token := range user.AuthTokens {
			fmt.Printf("Backend: Comparing token program '%s' with requested program '%s'\n", token.ProgramName, programName)
			if token.ProgramName == programName {
				isAuthorized = true
				authReason = fmt.Sprintf("Found matching auth token with program: %s", token.ProgramName)
				fmt.Printf("Backend: MATCH FOUND! Token program '%s' matches requested '%s'\n", token.ProgramName, programName)
				break
			}
		}

		if !isAuthorized {
			authReason = fmt.Sprintf("No matching program found. User has tokens for: %v", func() []string {
				programs := make([]string, len(user.AuthTokens))
				for i, token := range user.AuthTokens {
					programs[i] = token.ProgramName
				}
				return programs
			}())
		}
	}

	fmt.Printf("Backend: Authorization result: %v (Reason: %s)\n", isAuthorized, authReason)

	if !isAuthorized {
		fmt.Println("Backend: Final authorization check failed.")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorised to view this program.",
			"debug": authReason, // Remove this in production
		})
		return
	}
	fmt.Println("Backend: Final authorization check passed.")

	// 3. If authorized, proceed to fetch and return the workout data
	var program models.WorkoutProgram
	var dbProgramName string

	switch programName {
	case "beginner-program":
		dbProgramName = "beginner-program"
	case "advanced-program":
		dbProgramName = "advanced-program"
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
		return
	}

	if err := wc.DB.Where("name = ?", dbProgramName).First(&program).Error; err != nil {
		fmt.Printf("Backend: Program '%s' not found in database: %v\n", dbProgramName, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
		return
	}

	var workoutDay models.WorkoutDay
	if err := wc.DB.Where("program_id = ? AND day_number = ?", program.ID, dayNumber).
		Preload("WorkoutBlocks.Exercises.Exercise").
		First(&workoutDay).Error; err != nil {
		fmt.Printf("Backend: Workout day %d not found for program %d: %v\n", dayNumber, program.ID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout day not found for this program"})
		return
	}

	// log.Printf("Authorization check - UserID: %v, ProgramName: %s", userID, programName)
	// log.Printf("User role: %s, AuthTokens: %d", user.Role, len(user.AuthTokens))
	// for i, token := range user.AuthTokens {
	// 	log.Printf("  Token %d: Program=%s, Used=%v", i, token.ProgramName, token.IsUsed)
	// }
	// log.Printf("Is authorized: %v", isAuthorized)

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
