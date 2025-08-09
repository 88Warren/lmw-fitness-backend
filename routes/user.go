package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.POST("/api/register", uc.RegisterUser)
	router.POST("/api/login", uc.LoginUser)
	router.POST("/api/forgot-password", uc.RequestPasswordReset)
	router.POST("/api/verify-reset-token", uc.VerifyResetToken)
	router.POST("/api/reset-password", uc.ResetPassword)

	authenticated := router.Group("/api")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.PUT("/change-password-first-login", uc.ChangePassword)
		authenticated.GET("/profile", uc.GetProfile)
	}
}
