package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/utils/email"
	"github.com/88warren/lmw-fitness-backend/utils/emailtemplates"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsletterController struct {
	DB *gorm.DB
}

func NewNewsletterController(db *gorm.DB) *NewsletterController {
	return &NewsletterController{DB: db}
}

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (nc *NewsletterController) Subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingSubscriber models.NewsletterSubscriber
	if nc.DB.Where("email = ?", req.Email).First(&existingSubscriber).Error == nil {
		if existingSubscriber.IsConfirmed {
			ctx.JSON(http.StatusOK, gin.H{"message": "You are already subscribed to our newsletter!"})
			return
		}

		if err := nc.sendConfirmationEmail(existingSubscriber); err != nil {
			log.Printf("Error resending confirmation email: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resend confirmation email."})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "A confirmation email has been sent. Please check your inbox to confirm your subscription."})
		return
	}

	token, err := generateSecureToken(32)
	if err != nil {
		log.Printf("Error generating confirmation token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate confirmation token."})
		return
	}

	newSubscriber := models.NewsletterSubscriber{
		Email:        req.Email,
		SubscribedAt: time.Now(),
		IsConfirmed:  false,
		ConfirmToken: token,
	}

	if result := nc.DB.Create(&newSubscriber); result.Error != nil {
		log.Printf("Error creating newsletter subscriber: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe."})
		return
	}

	if err := nc.sendConfirmationEmail(newSubscriber); err != nil {
		log.Printf("Error sending confirmation email: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Subscription successful! Please check your inbox to confirm your subscription."})
}

func (nc *NewsletterController) ConfirmSubscription(ctx *gin.Context) {
	token := ctx.Param("token")

	var subscriber models.NewsletterSubscriber
	if nc.DB.Where("confirm_token = ?", token).First(&subscriber).Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired confirmation link."})
		return
	}

	if subscriber.IsConfirmed {
		ctx.JSON(http.StatusOK, gin.H{"message": "Your subscription is already confirmed!"})
		return
	}

	subscriber.IsConfirmed = true
	subscriber.ConfirmToken = ""
	if result := nc.DB.Save(&subscriber); result.Error != nil {
		log.Printf("Error confirming subscription: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm subscription."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully!"})
}

func (nc *NewsletterController) sendConfirmationEmail(subscriber models.NewsletterSubscriber) error {
	confirmLink := fmt.Sprintf("%s/newsletter/confirm/%s", os.Getenv("ALLOWED_ORIGIN"), subscriber.ConfirmToken)
	emailSubject := "Please confirm your LMW Fitness newsletter subscription"

	emailBody := emailtemplates.GenerateNewsletterConfirmationEmailBody(subscriber.Email, confirmLink)

	smtpPassword := getSMTPPasswordFromSecrets()

	return email.SendEmail(
		os.Getenv("SMTP_FROM"),
		subscriber.Email,
		emailSubject,
		emailBody,
		"",
		smtpPassword,
	)
}
