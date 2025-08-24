package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	api := router.Group("/api")
	{
		api.POST("/register", uc.RegisterUser)
		api.POST("/login", uc.LoginUser)
		api.POST("/forgot-password", uc.RequestPasswordReset)
		api.POST("/verify-reset-token", uc.VerifyResetToken)
		api.POST("/reset-password", uc.ResetPassword)
		api.POST("/verify-workout-token", uc.VerifyWorkoutToken)
	}
	authenticated := api.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.PUT("/change-password-first-login", uc.ChangePassword)
		authenticated.POST("/set-first-time-password", uc.SetFirstTimePassword)
		authenticated.GET("/profile", uc.GetProfile)
	}
}
