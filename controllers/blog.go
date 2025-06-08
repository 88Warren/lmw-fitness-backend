package controllers

import (
	"net/http"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BlogController struct {
	DB *gorm.DB
}

func NewBlogController(db *gorm.DB) *BlogController {
	return &BlogController{DB: db}
}

func (bc *BlogController) GetBlog(ctx *gin.Context) {
	var blogs []models.Blog
	if result := bc.DB.Find(&blogs); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blog posts: " + result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
}

func (bc *BlogController) GetBlogByID(ctx *gin.Context) {
	id := ctx.Param("id")

	var blog models.Blog
	if result := bc.DB.First(&blog, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog post: " + result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blog)
}

func (bc *BlogController) CreateBlog(ctx *gin.Context) {
	var newBlog models.Blog
	if err := ctx.ShouldBindJSON(&newBlog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := bc.DB.Create(&newBlog); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog post"})
		return
	}

	ctx.JSON(http.StatusCreated, newBlog)
}

func (bc *BlogController) UpdateBlog(ctx *gin.Context) {
	id := ctx.Param("id")

	var updatedBlog models.Blog
	if err := ctx.ShouldBindJSON(&updatedBlog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingBlog models.Blog
	if result := bc.DB.First(&existingBlog, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blog post: " + result.Error.Error()})
		return
	}

	if result := bc.DB.Model(&existingBlog).Updates(updatedBlog); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog post: " + result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, existingBlog)
}

func (bc *BlogController) DeleteBlog(ctx *gin.Context) {
	id := ctx.Param("id")

	if result := bc.DB.Delete(&models.Blog{}, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog post: " + result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
