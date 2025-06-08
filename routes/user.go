package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
	"github.com/laurawarren88/LMW_Fitness/middleware"
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
		authenticated.GET("/profile", uc.GetProfile)
	}
}
