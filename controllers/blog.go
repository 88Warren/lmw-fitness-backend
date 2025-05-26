package controllers

import (
	"net/http"

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

func (bc *BlogController) GetBlog(ctx *gin.Context) {
	var blogs []models.Blog
	if result := bc.DB.Find(&blogs); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blog posts: " + result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusOK, blogs)
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
