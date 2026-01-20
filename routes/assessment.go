package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAssessmentRoutes(router *gin.Engine, ac *controllers.AssessmentController) {
	// All assessment routes require authentication
	authenticated := router.Group("/api/assessments")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// Save or update assessment result
		authenticated.POST("/save", ac.SaveAssessment)

		// Get all assessment history for user
		authenticated.GET("/history", ac.GetAssessmentHistory)

		// Get comparison between Day 1 and Day 30 for a program
		authenticated.GET("/compare/:programName", ac.GetAssessmentComparison)

		// Get assessments for specific program and day
		authenticated.GET("/:programName/day/:dayNumber", ac.GetProgramAssessments)

		// Delete specific assessment
		authenticated.DELETE("/:id", ac.DeleteAssessment)
	}
}
