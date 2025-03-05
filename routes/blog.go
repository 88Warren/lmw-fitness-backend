package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterBlogRoutes(router *gin.Engine, bc *controllers.BlogController) {
	blogGroup := router.Group("/api/blogs")
	{
		blogGroup.GET("/", bc.GetBlogs)
		blogGroup.GET("/new", bc.CreateBlogForm)
		blogGroup.POST("/new", bc.CreateBlog)
		// blogGroup.GET("/:id", bc.GetBlogByID)
		// blogGroup.PUT("/:id", bc.UpdateBlog)
		// blogGroup.DELETE("/:id", bc.DeleteBlog)
	}
}
