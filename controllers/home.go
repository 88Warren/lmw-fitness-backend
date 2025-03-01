package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HomeController struct {
	DB *gorm.DB
}

func NewHomeController(db *gorm.DB) *HomeController {
	return &HomeController{DB: db}
}

type ContactForm struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Message string `json:"message" binding:"required"`
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Home Page",
	})
}

func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
	var form ContactForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process form submission (e.g., store in DB, send email)
	log.Printf("Contact form submitted: %+v\n", form)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message received!"})
}
