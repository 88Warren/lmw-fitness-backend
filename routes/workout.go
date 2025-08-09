package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterWorkoutRoutes(router *gin.Engine, wc *controllers.WorkoutController) {
	// Public routes for Browse workout programs.
	router.GET("/api/workouts/programs", wc.GetWorkoutPrograms)

	// Authenticated routes for managing user-specific workouts.
	authenticated := router.Group("/api/workouts")
	authenticated.Use(middleware.AuthMiddleware()) // Assumes you have an auth middleware
	{
		// GET retrieves data. Use URL params for specific resources.
		authenticated.GET("/programs/:programID/days/:dayNumber", wc.GetWorkoutDay)

		// POST creates a new resource or starts a new process.
		authenticated.POST("/start", wc.StartWorkout)
		authenticated.POST("/complete-exercise", wc.CompleteExercise)
		authenticated.POST("/complete-day", wc.CompleteWorkoutDay)

		// GET retrieves user-specific progress.
		authenticated.GET("/progress", wc.GetUserProgress)
	}
}
