package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
	"github.com/laurawarren88/LMW_Fitness/middleware"
)

func RegisterBlogRoutes(router *gin.Engine, bc *controllers.BlogController) {
	router.GET("/api/blog", bc.GetBlog)

	authenticated := router.Group("/api")
	authenticated.Use(middleware.AuthMiddleware())

	{
		authenticated.GET("/protected", func(c *gin.Context) {
			userID := c.MustGet("userID").(uint)
			userEmail := c.MustGet("userEmail").(string)
			userRole := c.MustGet("userRole").(string)
			c.JSON(http.StatusOK, gin.H{"message": "Welcome, authenticated user!", "userID": userID, "email": userEmail, "role": userRole})
		})
		authenticated.POST("/blog", middleware.RoleMiddleware("admin"), bc.CreateBlog)    // Requires admin role
		authenticated.PUT("/blog/:id", middleware.RoleMiddleware("admin"), bc.UpdateBlog) // Requires admin role
		authenticated.DELETE("/blog/:id", middleware.RoleMiddleware("admin"), bc.DeleteBlog)
	}
}
