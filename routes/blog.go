package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterBlogRoutes(router *gin.Engine, bc *controllers.BlogController) {
	// API routes
	api := router.Group("/api")
	{
		api.GET("/blog", bc.GetBlog)
		api.POST("/blog", bc.CreateBlog)
	}
}
