package controllers

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
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
	Token   string `json:"token"`
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Home Page",
	})
}

func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
	log.Println("Received contact form request")

	var form ContactForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		log.Printf("Form binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Form data received: %+v", form)
	log.Printf("reCAPTCHA token received: %s", form.Token)

	// if !verifyRecaptcha(form.Token) {
	// 	log.Printf("reCAPTCHA verification failed for token: %s", form.Token)
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
	// 	return
	// }

	err := sendEmail(form.Name, form.Email, form.Message)
	if err != nil {
		log.Printf("Email sending failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Process form submission (e.g., store in DB, send email)
	log.Printf("Contact form submitted: %+v\n", form)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message received!"})
}

func sendEmail(name, email, message string) error {

	log.Printf("SMTP Config - Host: %s, Port: %s, From: %s, To: %s",
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_TO"))

	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_FROM")) // Sender email from environment variable
	m.SetHeader("To", os.Getenv("SMTP_TO"))
	m.SetHeader("Reply-To", email) // Recipient email
	m.SetHeader("Subject", "New Contact Form Submission")

	m.SetBody("text/plain", "Name: "+name+"\nEmail: "+email+"\n\nMessage:\n"+message)

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	d.TLSConfig = &tls.Config{
		ServerName: os.Getenv("SMTP_HOST"),
	}

	return d.DialAndSend(m)
}

// func verifyRecaptcha(token string) bool {
// 	secret := os.Getenv("RECAPTCHA_SECRET")
// 	url := "https://www.google.com/recaptcha/api/siteverify"
// 	// Use form values instead of JSON for the request
// 	resp, err := http.PostForm(url, url.Values{
// 		"secret":   []string{secret},
// 		"response": []string{token},
// 	})
// 	if err != nil {
// 		log.Printf("reCAPTCHA verification error: %v", err)
// 		return false
// 	}
// 	defer resp.Body.Close()

// 	var result struct {
// 		Success bool `json:"success"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		log.Printf("reCAPTCHA response decode error: %v", err)
// 		return false
// 	}

// 	// Add debug logging
// 	log.Printf("reCAPTCHA verification result: %+v", result)
// 	return result.Success
// }
