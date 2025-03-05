package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterUserRoutes(router *gin.Engine, uc *controllers.UserController) {
	// userRoutes := router.Group("/api/users")
	// {
	// 	userRoutes.GET("/register", uc.GetSignupForm)
	// 	userRoutes.POST("/register", uc.SignupUser)
	// 	userRoutes.GET("/login", uc.GetLoginForm)
	// 	userRoutes.POST("/login", uc.LoginUser)
	// 	userRoutes.GET("/forgot_password", uc.ForgotPassword)
	// 	userRoutes.POST("/forgot_password", uc.ResetPassword)
	// 	userRoutes.GET("/profile/:id", uc.GetProfile)
	// 	userRoutes.POST("/logout", uc.LogoutUser)
	// }
}
