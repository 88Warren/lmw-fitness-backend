package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/models"
	"gorm.io/gorm"
)

type BlogController struct {
	DB *gorm.DB
}

func NewBlogController(db *gorm.DB) *BlogController {
	return &BlogController{DB: db}
}

// GET all blogs (only published for public, all for admin)
func (bc *BlogController) GetBlogs(ctx *gin.Context) {
	var blogs []models.Blog
	admin := ctx.GetBool("is_admin")

	query := bc.DB
	if !admin {
		query = query.Where("is_published = ?", true)
	}

	if err := query.Find(&blogs).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve blogs"})
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

// // GET a single blog by slug
// func (bc *BlogController) GetBlogBySlug(ctx *gin.Context) {
// 	slug := ctx.Param("slug")
// 	var blog models.Blog

// 	if err := bc.DB.Where("slug = ?", slug).First(&blog).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, blog)
// }

func (bc *BlogController) CreateBlogForm(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"title": "Create a New Blog Post",
	})
}

// POST a new blog (Admin only)
func (bc *BlogController) CreateBlog(ctx *gin.Context) {
	// if !ctx.GetBool("is_admin") {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	var blog models.Blog

	image, _ := ctx.FormFile("image_url")
	if image != nil {
		// Save the image
		imagePath := "./uploads/" + image.Filename // Adjust path as needed
		if err := ctx.SaveUploadedFile(image, imagePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		}
		blog.ImageURL = imagePath
	}

	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Set timestamp and slug
	blog.PublishedAt = time.Now()
	blog.Slug = generateSlug(blog.Title)

	if err := bc.DB.Create(&blog).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	ctx.JSON(http.StatusCreated, blog)
}

// // PUT update a blog (Admin only)
// func (bc *BlogController) UpdateBlog(ctx *gin.Context) {
// 	if !ctx.GetBool("is_admin") {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := ctx.Param("id")
// 	var blog models.Blog

// 	if err := bc.DB.First(&blog, id).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
// 		return
// 	}

// 	if err := ctx.ShouldBindJSON(&blog); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
// 		return
// 	}

// 	if err := bc.DB.Save(&blog).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, blog)
// }

// // DELETE a blog (Admin only)
// func (bc *BlogController) DeleteBlog(ctx *gin.Context) {
// 	if !ctx.GetBool("is_admin") {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := ctx.Param("id")
// 	if err := bc.DB.Delete(&models.Blog{}, id).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
// }

// Helper function to generate a slug
func generateSlug(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "-"))
}
