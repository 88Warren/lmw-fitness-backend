package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

// Exercise Management
func (ac *AdminController) GetAllExercises(c *gin.Context) {
	var exercises []models.Exercise
	if err := ac.DB.Preload("Modification").Preload("Modification2").Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercises"})
		return
	}
	c.JSON(http.StatusOK, exercises)
}

func (ac *AdminController) GetExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var exercise models.Exercise
	if err := ac.DB.Preload("Modification").Preload("Modification2").First(&exercise, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exercise"})
		return
	}
	c.JSON(http.StatusOK, exercise)
}

func (ac *AdminController) CreateExercise(c *gin.Context) {
	var exercise models.Exercise
	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Create(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

func (ac *AdminController) UpdateExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var exercise models.Exercise
	if err := ac.DB.First(&exercise, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find exercise"})
		return
	}

	if err := c.ShouldBindJSON(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Save(&exercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func (ac *AdminController) DeleteExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	if err := ac.DB.Delete(&models.Exercise{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}

// Workout Program Management
func (ac *AdminController) GetAllPrograms(c *gin.Context) {
	var programs []models.WorkoutProgram
	if err := ac.DB.Preload("Days").Find(&programs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve programs"})
		return
	}
	c.JSON(http.StatusOK, programs)
}

func (ac *AdminController) GetProgram(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	var program models.WorkoutProgram
	if err := ac.DB.Preload("Days.WorkoutBlocks.Exercises.Exercise").First(&program, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve program"})
		return
	}
	c.JSON(http.StatusOK, program)
}

func (ac *AdminController) CreateProgram(c *gin.Context) {
	var program models.WorkoutProgram
	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Create(&program).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create program"})
		return
	}

	c.JSON(http.StatusCreated, program)
}

func (ac *AdminController) UpdateProgram(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	var program models.WorkoutProgram
	if err := ac.DB.First(&program, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find program"})
		return
	}

	if err := c.ShouldBindJSON(&program); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Save(&program).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update program"})
		return
	}

	c.JSON(http.StatusOK, program)
}

func (ac *AdminController) DeleteProgram(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid program ID"})
		return
	}

	// Delete associated workout days and their blocks/exercises
	if err := ac.DB.Where("program_id = ?", id).Delete(&models.WorkoutDay{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated workout days"})
		return
	}

	if err := ac.DB.Delete(&models.WorkoutProgram{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete program"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Program deleted successfully"})
}

// Workout Day Management
func (ac *AdminController) CreateWorkoutDay(c *gin.Context) {
	var requestData struct {
		ProgramID     uint                  `json:"programId"`
		DayNumber     int                   `json:"dayNumber"`
		Title         string                `json:"title"`
		Description   string                `json:"description"`
		WorkoutBlocks []models.WorkoutBlock `json:"workoutBlocks"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := ac.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the workout day
	workoutDay := models.WorkoutDay{
		ProgramID:   requestData.ProgramID,
		DayNumber:   requestData.DayNumber,
		Title:       requestData.Title,
		Description: requestData.Description,
	}

	if err := tx.Create(&workoutDay).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout day"})
		return
	}

	// Create workout blocks and exercises
	for _, blockData := range requestData.WorkoutBlocks {
		block := models.WorkoutBlock{
			DayID:       workoutDay.ID,
			BlockType:   blockData.BlockType,
			BlockRounds: blockData.BlockRounds,
			RoundRest:   blockData.RoundRest,
			BlockNotes:  blockData.BlockNotes,
		}

		if err := tx.Create(&block).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout block"})
			return
		}

		// Create exercises for this block
		for _, exerciseData := range blockData.Exercises {
			exercise := models.WorkoutExercise{
				BlockID:       block.ID,
				ExerciseID:    exerciseData.ExerciseID,
				Order:         exerciseData.Order,
				Reps:          exerciseData.Reps,
				Duration:      exerciseData.Duration,
				WorkRestRatio: exerciseData.WorkRestRatio,
				Rest:          exerciseData.Rest,
				Tips:          exerciseData.Tips,
				Instructions:  exerciseData.Instructions,
			}

			if err := tx.Create(&exercise).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout exercise"})
				return
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	// Fetch the created workout day with all relations
	var createdWorkoutDay models.WorkoutDay
	if err := ac.DB.Where("id = ?", workoutDay.ID).
		Preload("WorkoutBlocks.Exercises.Exercise").
		First(&createdWorkoutDay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created workout day"})
		return
	}

	c.JSON(http.StatusCreated, createdWorkoutDay)
}

func (ac *AdminController) UpdateWorkoutDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout day ID"})
		return
	}

	var requestData struct {
		ProgramID     uint                  `json:"programId"`
		DayNumber     int                   `json:"dayNumber"`
		Title         string                `json:"title"`
		Description   string                `json:"description"`
		WorkoutBlocks []models.WorkoutBlock `json:"workoutBlocks"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := ac.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find the existing workout day
	var workoutDay models.WorkoutDay
	if err := tx.First(&workoutDay, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout day not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find workout day"})
		return
	}

	// Update basic workout day fields
	workoutDay.ProgramID = requestData.ProgramID
	workoutDay.DayNumber = requestData.DayNumber
	workoutDay.Title = requestData.Title
	workoutDay.Description = requestData.Description

	if err := tx.Save(&workoutDay).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workout day"})
		return
	}

	// Delete existing workout blocks and their exercises in the correct order
	// First, get all block IDs for this workout day
	var blockIDs []uint
	if err := tx.Model(&models.WorkoutBlock{}).Where("day_id = ?", workoutDay.ID).Pluck("id", &blockIDs).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing block IDs"})
		return
	}

	// Delete exercises first (they reference blocks)
	if len(blockIDs) > 0 {
		if err := tx.Where("block_id IN ?", blockIDs).Delete(&models.WorkoutExercise{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing exercises"})
			return
		}
	}

	// Then delete blocks (they reference the workout day)
	if err := tx.Where("day_id = ?", workoutDay.ID).Delete(&models.WorkoutBlock{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing blocks"})
		return
	}

	// Create new workout blocks and exercises
	for _, blockData := range requestData.WorkoutBlocks {
		block := models.WorkoutBlock{
			DayID:       workoutDay.ID,
			BlockType:   blockData.BlockType,
			BlockRounds: blockData.BlockRounds,
			RoundRest:   blockData.RoundRest,
			BlockNotes:  blockData.BlockNotes,
		}

		if err := tx.Create(&block).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout block"})
			return
		}

		// Create exercises for this block
		for _, exerciseData := range blockData.Exercises {
			exercise := models.WorkoutExercise{
				BlockID:       block.ID,
				ExerciseID:    exerciseData.ExerciseID,
				Order:         exerciseData.Order,
				Reps:          exerciseData.Reps,
				Duration:      exerciseData.Duration,
				WorkRestRatio: exerciseData.WorkRestRatio,
				Rest:          exerciseData.Rest,
				Tips:          exerciseData.Tips,
				Instructions:  exerciseData.Instructions,
			}

			if err := tx.Create(&exercise).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout exercise"})
				return
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
		return
	}

	// Fetch the updated workout day with all relations
	var updatedWorkoutDay models.WorkoutDay
	if err := ac.DB.Where("id = ?", workoutDay.ID).
		Preload("WorkoutBlocks.Exercises.Exercise").
		First(&updatedWorkoutDay).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated workout day"})
		return
	}

	c.JSON(http.StatusOK, updatedWorkoutDay)
}

func (ac *AdminController) DeleteWorkoutDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout day ID"})
		return
	}

	// Delete associated workout blocks and exercises
	var workoutBlocks []models.WorkoutBlock
	if err := ac.DB.Where("day_id = ?", id).Find(&workoutBlocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find workout blocks"})
		return
	}

	for _, block := range workoutBlocks {
		ac.DB.Where("block_id = ?", block.ID).Delete(&models.WorkoutExercise{})
	}
	ac.DB.Where("day_id = ?", id).Delete(&models.WorkoutBlock{})

	if err := ac.DB.Delete(&models.WorkoutDay{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workout day"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout day deleted successfully"})
}

// Workout Block Management
func (ac *AdminController) CreateWorkoutBlock(c *gin.Context) {
	var workoutBlock models.WorkoutBlock
	if err := c.ShouldBindJSON(&workoutBlock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Create(&workoutBlock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout block"})
		return
	}

	c.JSON(http.StatusCreated, workoutBlock)
}

func (ac *AdminController) UpdateWorkoutBlock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout block ID"})
		return
	}

	var workoutBlock models.WorkoutBlock
	if err := ac.DB.First(&workoutBlock, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout block not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find workout block"})
		return
	}

	if err := c.ShouldBindJSON(&workoutBlock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Save(&workoutBlock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workout block"})
		return
	}

	c.JSON(http.StatusOK, workoutBlock)
}

func (ac *AdminController) DeleteWorkoutBlock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout block ID"})
		return
	}

	// Delete associated workout exercises
	ac.DB.Where("block_id = ?", id).Delete(&models.WorkoutExercise{})

	if err := ac.DB.Delete(&models.WorkoutBlock{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workout block"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout block deleted successfully"})
}

// Workout Exercise Management
func (ac *AdminController) CreateWorkoutExercise(c *gin.Context) {
	var workoutExercise models.WorkoutExercise
	if err := c.ShouldBindJSON(&workoutExercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Create(&workoutExercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workout exercise"})
		return
	}

	c.JSON(http.StatusCreated, workoutExercise)
}

func (ac *AdminController) UpdateWorkoutExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout exercise ID"})
		return
	}

	var workoutExercise models.WorkoutExercise
	if err := ac.DB.First(&workoutExercise, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout exercise not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find workout exercise"})
		return
	}

	if err := c.ShouldBindJSON(&workoutExercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ac.DB.Save(&workoutExercise).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workout exercise"})
		return
	}

	c.JSON(http.StatusOK, workoutExercise)
}

func (ac *AdminController) DeleteWorkoutExercise(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workout exercise ID"})
		return
	}

	if err := ac.DB.Delete(&models.WorkoutExercise{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workout exercise"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workout exercise deleted successfully"})
}

// User Management
func (ac *AdminController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := ac.DB.Preload("AuthTokens").Preload("UserPrograms.WorkoutProgram").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	// Enhanced user data for admin monitoring
	var enhancedUsers []map[string]interface{}
	for _, user := range users {
		// Calculate purchased programs from UserPrograms (consistent with user profile)
		purchasedPrograms := make(map[string]bool)
		for _, userProgram := range user.UserPrograms {
			if userProgram.WorkoutProgram.Name != "" {
				purchasedPrograms[userProgram.WorkoutProgram.Name] = true
			}
		}
		programList := make([]string, 0, len(purchasedPrograms))
		for program := range purchasedPrograms {
			programList = append(programList, program)
		}

		// Calculate workout progress
		totalCompletedDays := 0
		for _, count := range user.CompletedDays {
			totalCompletedDays += count
		}

		// Calculate account age
		accountAge := time.Since(user.CreatedAt).Hours() / 24

		// Determine user activity level
		activityLevel := "New"
		if totalCompletedDays > 0 {
			activityLevel = "Active"
		}
		if totalCompletedDays > 10 {
			activityLevel = "Regular"
		}
		if totalCompletedDays > 30 {
			activityLevel = "Dedicated"
		}

		// Calculate last activity (approximation based on updated_at)
		daysSinceLastActivity := time.Since(user.UpdatedAt).Hours() / 24

		enhancedUser := map[string]interface{}{
			"ID":                    user.ID,
			"email":                 user.Email,
			"role":                  user.Role,
			"isActive":              true, // Default to true since we don't have this field
			"CreatedAt":             user.CreatedAt,
			"UpdatedAt":             user.UpdatedAt,
			"lastLogin":             user.UpdatedAt, // Approximation
			"purchasedPrograms":     programList,
			"totalCompletedDays":    totalCompletedDays,
			"accountAge":            int(accountAge),
			"activityLevel":         activityLevel,
			"daysSinceLastActivity": int(daysSinceLastActivity),
			"timezone":              user.Timezone,
			"mustChangePassword":    user.MustChangePassword,
			"programCount":          len(programList),
			"authTokenCount":        len(user.AuthTokens),
		}
		enhancedUsers = append(enhancedUsers, enhancedUser)
	}

	c.JSON(http.StatusOK, enhancedUsers)
}

func (ac *AdminController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := ac.DB.Select("id, email, role, created_at, updated_at").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (ac *AdminController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := ac.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	var updateData struct {
		Role     string `json:"role"`
		IsActive *bool  `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the fields that are provided
	updates := make(map[string]interface{})
	if updateData.Role != "" {
		if updateData.Role != "user" && updateData.Role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Must be 'user' or 'admin'"})
			return
		}
		updates["role"] = updateData.Role
	}
	if updateData.IsActive != nil {
		updates["is_active"] = *updateData.IsActive
	}

	if err := ac.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return updated user data (excluding sensitive fields)
	var updatedUser models.User
	if err := ac.DB.Select("id, email, role, created_at, updated_at").First(&updatedUser, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (ac *AdminController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	var user models.User
	if err := ac.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Prevent deletion of the last admin user
	if user.Role == "admin" {
		var adminCount int64
		ac.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the last admin user"})
			return
		}
	}

	if err := ac.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
func (ac *AdminController) ResetUserPassword(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user exists
	var user models.User
	if err := ac.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	// Reuse the password reset logic by creating a UserController instance
	userController := NewUserController(ac.DB)

	// Create a mock request body with the user's email
	mockBody := fmt.Sprintf(`{"email":"%s"}`, user.Email)
	req, _ := http.NewRequest("POST", "/forgot-password", strings.NewReader(mockBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder and mock context
	w := httptest.NewRecorder()
	mockCtx, _ := gin.CreateTestContext(w)
	mockCtx.Request = req

	// Call the existing RequestPasswordReset method
	userController.RequestPasswordReset(mockCtx)

	// Check if the request was successful
	if w.Code != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent successfully",
		"email":   user.Email,
	})
}

// Analytics Dashboard
func (ac *AdminController) GetAnalyticsDashboard(c *gin.Context) {
	// Get time range from query params (default to last 30 days)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// User Analytics
	var totalUsers int64
	var newUsers int64
	var activeUsers int64
	var adminUsers int64

	ac.DB.Model(&models.User{}).Count(&totalUsers)
	ac.DB.Model(&models.User{}).Where("created_at >= ?", startDate).Count(&newUsers)
	ac.DB.Model(&models.User{}).Where("updated_at >= ?", startDate).Count(&activeUsers)
	ac.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminUsers)

	// Program Analytics
	var totalPrograms int64
	var activePrograms int64
	var totalWorkoutDays int64
	var totalExercises int64

	ac.DB.Model(&models.WorkoutProgram{}).Count(&totalPrograms)
	ac.DB.Model(&models.WorkoutProgram{}).Where("is_active = ?", true).Count(&activePrograms)
	ac.DB.Model(&models.WorkoutDay{}).Count(&totalWorkoutDays)
	ac.DB.Model(&models.Exercise{}).Count(&totalExercises)

	// Content Analytics
	var totalBlogs int64
	var featuredBlogs int64
	var newsletterSubscribers int64
	var confirmedSubscribers int64

	ac.DB.Model(&models.Blog{}).Count(&totalBlogs)
	ac.DB.Model(&models.Blog{}).Where("is_featured = ?", true).Count(&featuredBlogs)
	ac.DB.Model(&models.NewsletterSubscriber{}).Count(&newsletterSubscribers)
	ac.DB.Model(&models.NewsletterSubscriber{}).Where("is_confirmed = ?", true).Count(&confirmedSubscribers)

	// User Registration Trends (last 7 days)
	var registrationTrends []map[string]interface{}
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var count int64
		ac.DB.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).Count(&count)

		registrationTrends = append(registrationTrends, map[string]interface{}{
			"date":  startOfDay.Format("2006-01-02"),
			"count": count,
		})
	}

	// Program Popularity
	var programStats []map[string]interface{}
	var programs []models.WorkoutProgram
	ac.DB.Preload("Days").Find(&programs)

	for _, program := range programs {
		var userCount int64
		ac.DB.Model(&models.UserProgram{}).Where("program_id = ?", program.ID).Count(&userCount)

		programStats = append(programStats, map[string]interface{}{
			"name":       program.Name,
			"difficulty": program.Difficulty,
			"duration":   program.Duration,
			"userCount":  userCount,
			"dayCount":   len(program.Days),
			"isActive":   program.IsActive,
		})
	}

	// User Activity Levels
	var users []models.User
	ac.DB.Find(&users)

	activityLevels := map[string]int{
		"New":       0,
		"Active":    0,
		"Regular":   0,
		"Dedicated": 0,
	}

	for _, user := range users {
		totalCompletedDays := 0
		for _, count := range user.CompletedDays {
			totalCompletedDays += count
		}

		if totalCompletedDays == 0 {
			activityLevels["New"]++
		} else if totalCompletedDays <= 10 {
			activityLevels["Active"]++
		} else if totalCompletedDays <= 30 {
			activityLevels["Regular"]++
		} else {
			activityLevels["Dedicated"]++
		}
	}

	// Exercise Category Distribution
	var exerciseCategories []map[string]interface{}
	rows, err := ac.DB.Model(&models.Exercise{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var category string
			var count int64
			rows.Scan(&category, &count)
			exerciseCategories = append(exerciseCategories, map[string]interface{}{
				"category": category,
				"count":    count,
			})
		}
	}

	// Recent Activity (last 10 users)
	var recentUsers []models.User
	ac.DB.Order("created_at DESC").Limit(10).Find(&recentUsers)

	var recentActivity []map[string]interface{}
	for _, user := range recentUsers {
		recentActivity = append(recentActivity, map[string]interface{}{
			"email":     user.Email,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
			"timezone":  user.Timezone,
		})
	}

	// System Health Metrics
	systemHealth := map[string]interface{}{
		"totalUsers":            totalUsers,
		"totalPrograms":         totalPrograms,
		"totalWorkoutDays":      totalWorkoutDays,
		"totalExercises":        totalExercises,
		"totalBlogs":            totalBlogs,
		"newsletterSubscribers": newsletterSubscribers,
		"databaseConnected":     true, // Since we're querying successfully
		"lastUpdated":           time.Now(),
	}

	analytics := map[string]interface{}{
		"overview": map[string]interface{}{
			"totalUsers":            totalUsers,
			"newUsers":              newUsers,
			"activeUsers":           activeUsers,
			"adminUsers":            adminUsers,
			"totalPrograms":         totalPrograms,
			"activePrograms":        activePrograms,
			"totalWorkoutDays":      totalWorkoutDays,
			"totalExercises":        totalExercises,
			"totalBlogs":            totalBlogs,
			"featuredBlogs":         featuredBlogs,
			"newsletterSubscribers": newsletterSubscribers,
			"confirmedSubscribers":  confirmedSubscribers,
		},
		"trends": map[string]interface{}{
			"registrationTrends": registrationTrends,
			"timeRange":          fmt.Sprintf("Last %d days", days),
		},
		"programs": map[string]interface{}{
			"programStats":   programStats,
			"activityLevels": activityLevels,
		},
		"content": map[string]interface{}{
			"exerciseCategories": exerciseCategories,
			"recentActivity":     recentActivity,
		},
		"system":      systemHealth,
		"generatedAt": time.Now(),
	}

	c.JSON(http.StatusOK, analytics)
}
