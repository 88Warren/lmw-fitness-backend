package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterWorkoutRoutes(router *gin.Engine, wc *controllers.WorkoutController) {
	router.GET("/api/workouts/programs", wc.GetWorkoutPrograms)

	authenticated := router.Group("/api/workouts")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/programs/:programID/days/:dayNumber", wc.GetWorkoutDay)
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
