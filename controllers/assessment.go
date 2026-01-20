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

type AssessmentController struct {
	DB *gorm.DB
}

func NewAssessmentController(db *gorm.DB) *AssessmentController {
	return &AssessmentController{DB: db}
}

type SaveAssessmentRequest struct {
	ProgramName  string `json:"programName" binding:"required"`
	DayNumber    int    `json:"dayNumber" binding:"required"`
	ExerciseID   uint   `json:"exerciseId" binding:"required"`
	ExerciseName string `json:"exerciseName" binding:"required"`
	Reps         *int   `json:"reps"`
	TimeSeconds  *int   `json:"timeSeconds"`
	Notes        string `json:"notes"`
}

type AssessmentResponse struct {
	ID           uint      `json:"id"`
	ProgramName  string    `json:"programName"`
	DayNumber    int       `json:"dayNumber"`
	ExerciseID   uint      `json:"exerciseId"`
	ExerciseName string    `json:"exerciseName"`
	Reps         *int      `json:"reps"`
	TimeSeconds  *int      `json:"timeSeconds"`
	Notes        string    `json:"notes"`
	RecordedDate time.Time `json:"recordedDate"`
}

type ComparisonResponse struct {
	ExerciseName string              `json:"exerciseName"`
	Day1         *AssessmentResponse `json:"day1"`
	Day30        *AssessmentResponse `json:"day30"`
	Improvement  *ImprovementMetrics `json:"improvement"`
}

type ImprovementMetrics struct {
	RepsDifference  *int     `json:"repsDifference"`
	TimeDifference  *int     `json:"timeDifference"`
	PercentImproved *float64 `json:"percentImproved"`
}

// SaveAssessment saves a fitness assessment result
func (ac *AssessmentController) SaveAssessment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req SaveAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Log the binding error for debugging
		fmt.Printf("Assessment save binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Log the received request for debugging
	fmt.Printf("Assessment save request: %+v\n", req)

	// Validate that either reps or timeSeconds is provided
	if req.Reps == nil && req.TimeSeconds == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either reps or timeSeconds must be provided"})
		return
	}

	// Check if assessment already exists for this user, program, day, and exercise
	var existingAssessment models.FitnessAssessment
	result := ac.DB.Where("user_id = ? AND program_name = ? AND day_number = ? AND exercise_id = ?",
		userID, req.ProgramName, req.DayNumber, req.ExerciseID).First(&existingAssessment)

	assessment := models.FitnessAssessment{
		UserID:       userID.(uint),
		ProgramName:  req.ProgramName,
		DayNumber:    req.DayNumber,
		ExerciseID:   req.ExerciseID,
		ExerciseName: req.ExerciseName,
		Reps:         req.Reps,
		TimeSeconds:  req.TimeSeconds,
		Notes:        req.Notes,
		RecordedDate: time.Now(),
	}

	if result.Error == nil {
		// Update existing assessment
		assessment.ID = existingAssessment.ID
		if err := ac.DB.Save(&assessment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assessment"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Assessment updated successfully", "assessment": assessment})
	} else {
		// Create new assessment
		if err := ac.DB.Create(&assessment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save assessment"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Assessment saved successfully", "assessment": assessment})
	}
}

// GetAssessmentHistory gets all assessments for a user
func (ac *AssessmentController) GetAssessmentHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var assessments []models.FitnessAssessment
	if err := ac.DB.Where("user_id = ?", userID).
		Preload("Exercise").
		Order("recorded_date DESC").
		Find(&assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assessment history"})
		return
	}

	// Convert to response format
	var response []AssessmentResponse
	for _, assessment := range assessments {
		response = append(response, AssessmentResponse{
			ID:           assessment.ID,
			ProgramName:  assessment.ProgramName,
			DayNumber:    assessment.DayNumber,
			ExerciseID:   assessment.ExerciseID,
			ExerciseName: assessment.ExerciseName,
			Reps:         assessment.Reps,
			TimeSeconds:  assessment.TimeSeconds,
			Notes:        assessment.Notes,
			RecordedDate: assessment.RecordedDate,
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetAssessmentComparison compares Day 1 vs Day 30 assessments for a program
func (ac *AssessmentController) GetAssessmentComparison(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	programName := c.Param("programName")
	if programName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Program name is required"})
		return
	}

	// Get Day 1 assessments
	var day1Assessments []models.FitnessAssessment
	if err := ac.DB.Where("user_id = ? AND program_name = ? AND day_number = ?",
		userID, programName, 1).Find(&day1Assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Day 1 assessments"})
		return
	}

	// Get Day 30 assessments
	var day30Assessments []models.FitnessAssessment
	if err := ac.DB.Where("user_id = ? AND program_name = ? AND day_number = ?",
		userID, programName, 30).Find(&day30Assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Day 30 assessments"})
		return
	}

	// Create comparison map
	day30Map := make(map[uint]*models.FitnessAssessment)
	for i := range day30Assessments {
		day30Map[day30Assessments[i].ExerciseID] = &day30Assessments[i]
	}

	var comparisons []ComparisonResponse
	for _, day1 := range day1Assessments {
		comparison := ComparisonResponse{
			ExerciseName: day1.ExerciseName,
			Day1: &AssessmentResponse{
				ID:           day1.ID,
				ProgramName:  day1.ProgramName,
				DayNumber:    day1.DayNumber,
				ExerciseID:   day1.ExerciseID,
				ExerciseName: day1.ExerciseName,
				Reps:         day1.Reps,
				TimeSeconds:  day1.TimeSeconds,
				Notes:        day1.Notes,
				RecordedDate: day1.RecordedDate,
			},
		}

		if day30, exists := day30Map[day1.ExerciseID]; exists {
			comparison.Day30 = &AssessmentResponse{
				ID:           day30.ID,
				ProgramName:  day30.ProgramName,
				DayNumber:    day30.DayNumber,
				ExerciseID:   day30.ExerciseID,
				ExerciseName: day30.ExerciseName,
				Reps:         day30.Reps,
				TimeSeconds:  day30.TimeSeconds,
				Notes:        day30.Notes,
				RecordedDate: day30.RecordedDate,
			}

			// Calculate improvement metrics
			improvement := &ImprovementMetrics{}
			if day1.Reps != nil && day30.Reps != nil {
				diff := *day30.Reps - *day1.Reps
				improvement.RepsDifference = &diff
				if *day1.Reps > 0 {
					percent := float64(diff) / float64(*day1.Reps) * 100
					improvement.PercentImproved = &percent
				}
			}
			if day1.TimeSeconds != nil && day30.TimeSeconds != nil {
				diff := *day30.TimeSeconds - *day1.TimeSeconds
				improvement.TimeDifference = &diff
				if *day1.TimeSeconds > 0 {
					percent := float64(diff) / float64(*day1.TimeSeconds) * 100
					improvement.PercentImproved = &percent
				}
			}
			comparison.Improvement = improvement
		}

		comparisons = append(comparisons, comparison)
	}

	c.JSON(http.StatusOK, comparisons)
}

// GetProgramAssessments gets all assessments for a specific program and day
func (ac *AssessmentController) GetProgramAssessments(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	programName := c.Param("programName")
	dayNumberStr := c.Param("dayNumber")

	dayNumber, err := strconv.Atoi(dayNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day number"})
		return
	}

	var assessments []models.FitnessAssessment
	if err := ac.DB.Where("user_id = ? AND program_name = ? AND day_number = ?",
		userID, programName, dayNumber).
		Preload("Exercise").
		Order("exercise_id").
		Find(&assessments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve assessments"})
		return
	}

	// Convert to response format
	var response []AssessmentResponse
	for _, assessment := range assessments {
		response = append(response, AssessmentResponse{
			ID:           assessment.ID,
			ProgramName:  assessment.ProgramName,
			DayNumber:    assessment.DayNumber,
			ExerciseID:   assessment.ExerciseID,
			ExerciseName: assessment.ExerciseName,
			Reps:         assessment.Reps,
			TimeSeconds:  assessment.TimeSeconds,
			Notes:        assessment.Notes,
			RecordedDate: assessment.RecordedDate,
		})
	}

	c.JSON(http.StatusOK, response)
}

// DeleteAssessment deletes a specific assessment
func (ac *AssessmentController) DeleteAssessment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	assessmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assessment ID"})
		return
	}

	// Verify the assessment belongs to the user
	var assessment models.FitnessAssessment
	if err := ac.DB.Where("id = ? AND user_id = ?", assessmentID, userID).First(&assessment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Assessment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find assessment"})
		return
	}

	if err := ac.DB.Delete(&assessment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assessment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assessment deleted successfully"})
}
