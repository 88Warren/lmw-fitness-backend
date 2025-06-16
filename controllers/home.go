package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/88warren/lmw-fitness-backend/utils/email"
	"github.com/88warren/lmw-fitness-backend/utils/emailtemplates"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
	Token   string `json:"token"`
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Home Page",
	})
}

func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
	var form ContactForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		log.Printf("Form binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"stattus": "error",
		})
		return
	}

	_, smtpPassword := getK8sSecrets()

	log.Printf("Form data received: %+v", form)
	log.Printf("reCAPTCHA token received: %s", form.Token)

	if !verifyRecaptcha(form.Token) {
		log.Printf("reCAPTCHA verification failed for token: %s", form.Token)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
		return
	}

	emailBody := emailtemplates.GenerateContactFormEmailBody(form.Name, form.Email, form.Subject, form.Message)

	err := email.SendEmail(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_TO"),
		"New Contact Form Submission: "+form.Subject,
		emailBody,
		form.Email,
		smtpPassword,
	)
	if err != nil {
		log.Printf("Email sending failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message received!"})
}

func getK8sSecrets() (string, string) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return os.Getenv("RECAPTCHA_SECRET"), os.Getenv("SMTP_PASSWORD")
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to load cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	ctx := context.TODO()
	secret, err := clientset.CoreV1().Secrets("lmw-fitness").Get(ctx, "lmw-fitness-api-secrets", metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	recaptchaSecret := string(secret.Data["RECAPTCHA_SECRET"])
	smtpPassword := string(secret.Data["SMTP_PASSWORD"])

	return recaptchaSecret, smtpPassword
}

func verifyRecaptcha(token string) bool {
	secret := os.Getenv("RECAPTCHA_SECRET")
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	log.Printf("Verifying reCAPTCHA with secret: %s, token: %s", secret, token)
	log.Printf("Using RECAPTCHA_SECRET: %s", secret)

	resp, err := http.PostForm(verifyURL, url.Values{
		"secret":   {secret},
		"response": {token},
	})
	if err != nil {
		log.Printf("reCAPTCHA verification error: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success     bool     `json:"success"`
		ChallengeTs string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		ErrorCodes  []string `json:"error-codes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("reCAPTCHA response decode error: %v", err)
		return false
	}

	log.Printf("reCAPTCHA verification result: %+v", result)
	log.Printf("Using RECAPTCHA_SECRET: %s", secret)
	return result.Success
}
