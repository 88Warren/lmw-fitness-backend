package routes

import (
	"net/http"

	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterBlogRoutes(router *gin.Engine, bc *controllers.BlogController) {
	router.GET("/api/blog", bc.GetBlog)
	router.GET("/api/blog/:id", bc.GetBlogByID)

	authenticated := router.Group("/api")
	authenticated.Use(middleware.AuthMiddleware())

	{
		authenticated.GET("/protected", func(c *gin.Context) {
			userID := c.MustGet("userID").(uint)
			userEmail := c.MustGet("userEmail").(string)
			userRole := c.MustGet("userRole").(string)
			c.JSON(http.StatusOK, gin.H{"message": "Welcome, authenticated user!", "userID": userID, "email": userEmail, "role": userRole})
		})
		authenticated.POST("/blog", middleware.RoleMiddleware("admin"), bc.CreateBlog)
		authenticated.PUT("/blog/:id", middleware.RoleMiddleware("admin"), bc.UpdateBlog)
		authenticated.DELETE("/blog/:id", middleware.RoleMiddleware("admin"), bc.DeleteBlog)
	}
}
