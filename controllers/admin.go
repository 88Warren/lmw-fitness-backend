package controllers

import (
	"net/http"
	"strconv"

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
