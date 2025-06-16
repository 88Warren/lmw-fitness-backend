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
	// First try environment variables
	recaptchaSecret := os.Getenv("RECAPTCHA_SECRET")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	log.Printf("Initial environment variables - RECAPTCHA_SECRET: %s, SMTP_PASSWORD: %s",
		func() string {
			if len(recaptchaSecret) > 5 {
				return recaptchaSecret[:5] + "..."
			}
			return "not set"
		}(),
		func() string {
			if len(smtpPassword) > 5 {
				return smtpPassword[:5] + "..."
			}
			return "not set"
		}())

	// If we're in Kubernetes, try to get secrets from the cluster
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Printf("Failed to load cluster config: %v", err)
			return recaptchaSecret, smtpPassword
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			log.Printf("Failed to create Kubernetes client: %v", err)
			return recaptchaSecret, smtpPassword
		}

		ctx := context.TODO()
		secretName := os.Getenv("SECRET_NAME")
		if secretName == "" {
			secretName = "lmw-fitness-api-secrets"
		}

		log.Printf("Attempting to get secret: %s", secretName)

		secret, err := clientset.CoreV1().Secrets("lmw-fitness").Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get secret: %v", err)
			return recaptchaSecret, smtpPassword
		}

		if recaptchaSecret == "" {
			recaptchaSecret = string(secret.Data["RECAPTCHA_SECRET"])
		}
		if smtpPassword == "" {
			smtpPassword = string(secret.Data["SMTP_PASSWORD"])
		}

		log.Printf("Secrets from Kubernetes - RECAPTCHA_SECRET: %s, SMTP_PASSWORD: %s",
			func() string {
				if len(recaptchaSecret) > 5 {
					return recaptchaSecret[:5] + "..."
				}
				return "not set"
			}(),
			func() string {
				if len(smtpPassword) > 5 {
					return smtpPassword[:5] + "..."
				}
				return "not set"
			}())
	}

	return recaptchaSecret, smtpPassword
}

func verifyRecaptcha(token string) bool {
	secret, _ := getK8sSecrets()
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	log.Printf("Verifying reCAPTCHA with token: %s", token)
	log.Printf("Using RECAPTCHA_SECRET: %s", func() string {
		if len(secret) > 5 {
			return secret[:5] + "..."
		}
		return "not set"
	}())

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
	return result.Success
}
