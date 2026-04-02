package controllers

import (
	"net/http"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AMRAPController struct {
	DB *gorm.DB
}

func NewAMRAPController(db *gorm.DB) *AMRAPController {
	return &AMRAPController{DB: db}
}

type SaveAMRAPScoreRequest struct {
	BlockID     uint   `json:"blockId"`
	ProgramName string `json:"programName" binding:"required"`
	DayNumber   int    `json:"dayNumber" binding:"required"`
	BlockIndex  int    `json:"blockIndex"`
	Rounds      int    `json:"rounds"`
	PartialReps int    `json:"partialReps"`
	Notes       string `json:"notes"`
}

type AMRAPScoreResponse struct {
	ID           uint      `json:"id"`
	BlockID      uint      `json:"blockId"`
	ProgramName  string    `json:"programName"`
	DayNumber    int       `json:"dayNumber"`
	BlockIndex   int       `json:"blockIndex"`
	Rounds       int       `json:"rounds"`
	PartialReps  int       `json:"partialReps"`
	Notes        string    `json:"notes"`
	RecordedDate time.Time `json:"recordedDate"`
}

// SaveAMRAPScore saves or updates the user's best AMRAP score for a block
func (ac *AMRAPController) SaveAMRAPScore(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req SaveAMRAPScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BlockID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "blockId is required"})
		return
	}
	if req.Rounds < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rounds must be 0 or more"})
		return
	}

	// Check if a score already exists for this user + block
	var existing models.AMRAPScore
	result := ac.DB.Where("user_id = ? AND block_id = ?", userID, req.BlockID).First(&existing)

	score := models.AMRAPScore{
		UserID:       userID.(uint),
		BlockID:      req.BlockID,
		ProgramName:  req.ProgramName,
		DayNumber:    req.DayNumber,
		BlockIndex:   req.BlockIndex,
		Rounds:       req.Rounds,
		PartialReps:  req.PartialReps,
		Notes:        req.Notes,
		RecordedDate: time.Now(),
	}

	if result.Error == nil {
		// Only update if the new score is better (more rounds, or same rounds with more partial reps)
		isBetter := req.Rounds > existing.Rounds ||
			(req.Rounds == existing.Rounds && req.PartialReps > existing.PartialReps)

		if isBetter {
			score.ID = existing.ID
			if err := ac.DB.Save(&score).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update score"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message":   "New personal best saved!",
				"score":     toAMRAPResponse(score),
				"isNewBest": true,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message":   "Score recorded but didn't beat your personal best",
				"score":     toAMRAPResponse(existing),
				"isNewBest": false,
			})
		}
	} else {
		// First time recording this block
		if err := ac.DB.Create(&score).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save score"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message":   "Score saved!",
			"score":     toAMRAPResponse(score),
			"isNewBest": true,
		})
	}
}

// GetAMRAPScore gets the user's best score for a specific block
func (ac *AMRAPController) GetAMRAPScore(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	blockID := c.Param("blockId")

	var score models.AMRAPScore
	if err := ac.DB.Where("user_id = ? AND block_id = ?", userID, blockID).First(&score).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No score found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve score"})
		return
	}

	c.JSON(http.StatusOK, toAMRAPResponse(score))
}

// GetAllAMRAPScores gets all AMRAP scores for the user
func (ac *AMRAPController) GetAllAMRAPScores(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var scores []models.AMRAPScore
	if err := ac.DB.Where("user_id = ?", userID).
		Order("recorded_date DESC").
		Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve scores"})
		return
	}

	var response []AMRAPScoreResponse
	for _, s := range scores {
		response = append(response, toAMRAPResponse(s))
	}

	c.JSON(http.StatusOK, response)
}

func toAMRAPResponse(s models.AMRAPScore) AMRAPScoreResponse {
	return AMRAPScoreResponse{
		ID:           s.ID,
		BlockID:      s.BlockID,
		ProgramName:  s.ProgramName,
		DayNumber:    s.DayNumber,
		BlockIndex:   s.BlockIndex,
		Rounds:       s.Rounds,
		PartialReps:  s.PartialReps,
		Notes:        s.Notes,
		RecordedDate: s.RecordedDate,
	}
}
