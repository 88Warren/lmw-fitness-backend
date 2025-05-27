package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
	"github.com/laurawarren88/LMW_Fitness/middleware"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	router.POST("/api/register", uc.RegisterUser)
	router.POST("/api/login", uc.LoginUser)

	authenticated := router.Group("/api")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/profile", uc.GetProfile)
	}
}
