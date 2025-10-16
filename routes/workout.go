package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterWorkoutRoutes(router *gin.Engine, wc *controllers.WorkoutController) {
	// Public routes
	router.GET("/api/workouts/programs", wc.GetWorkoutPrograms)
	router.GET("/api/workouts/programs/:id", wc.GetWorkoutProgramByID)

	// Authenticated routes
	authenticated := router.Group("/api/workouts")
	authenticated.Use(middleware.AuthMiddleware())
	{
		// Use 'program' instead of 'programs' to avoid conflicts
		authenticated.GET("/program/:programID/days/:dayNumber", wc.GetWorkoutDay)
		authenticated.GET("/:programName/list", wc.GetProgramList)
		authenticated.GET("/:programName/routines/warmup", wc.GetWarmup)
		authenticated.GET("/:programName/routines/cooldown", wc.GetCooldown)
		authenticated.GET("/:programName/day/:dayNumber", wc.GetWorkoutDayByProgramAndDay)

		authenticated.POST("/start", wc.StartWorkout)
		authenticated.POST("/complete-exercise", wc.CompleteExercise)
		authenticated.POST("/complete-day", wc.CompleteWorkoutDay)

		authenticated.GET("/progress", wc.GetUserProgress)
	}
}
