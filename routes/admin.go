package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(router *gin.Engine, ac *controllers.AdminController) {
	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware()) // We'll need to create this middleware
	{
		// Exercise management
		admin.GET("/exercises", ac.GetAllExercises)
		admin.GET("/exercises/:id", ac.GetExercise)
		admin.POST("/exercises", ac.CreateExercise)
		admin.PUT("/exercises/:id", ac.UpdateExercise)
		admin.DELETE("/exercises/:id", ac.DeleteExercise)

		// Program management
		admin.GET("/programs", ac.GetAllPrograms)
		admin.GET("/programs/:id", ac.GetProgram)
		admin.POST("/programs", ac.CreateProgram)
		admin.PUT("/programs/:id", ac.UpdateProgram)
		admin.DELETE("/programs/:id", ac.DeleteProgram)

		// Workout day management
		admin.POST("/workout-days", ac.CreateWorkoutDay)
		admin.PUT("/workout-days/:id", ac.UpdateWorkoutDay)
		admin.DELETE("/workout-days/:id", ac.DeleteWorkoutDay)

		// Workout block management
		admin.POST("/workout-blocks", ac.CreateWorkoutBlock)
		admin.PUT("/workout-blocks/:id", ac.UpdateWorkoutBlock)
		admin.DELETE("/workout-blocks/:id", ac.DeleteWorkoutBlock)

		// Workout exercise management
		admin.POST("/workout-exercises", ac.CreateWorkoutExercise)
		admin.PUT("/workout-exercises/:id", ac.UpdateWorkoutExercise)
		admin.DELETE("/workout-exercises/:id", ac.DeleteWorkoutExercise)
	}
}
