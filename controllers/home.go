package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/utils/email"
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

// var allowedOrigins []string

// func init() {
// 	envAllowedOrigins := os.Getenv("ALLOWED_ORIGIN")
// 	if envAllowedOrigins == "" {
// 		// log.Println("WARNING: ALLOWED_ORIGIN environment variable is not set. CORS might be misconfigured.")
// 		allowedOrigins = []string{}
// 	} else {
// 		allowedOrigins = strings.Split(envAllowedOrigins, ",")
// 		// log.Printf("CORS allowed origins: %v", allowedOrigins)
// 	}
// }

// func CheckOrigin(origin string) string {
// 	for _, allowed := range allowedOrigins {
// 		if origin == allowed {
// 			return allowed
// 		}
// 	}
// 	return ""
// }

func (hc *HomeController) GetHome(ctx *gin.Context) {
	// origin := ctx.GetHeader("Origin")
	// if allowed := CheckOrigin(origin); allowed != "" {
	// 	ctx.Header("Access-Control-Allow-Origin", allowed)
	// }

	// ctx.Header("Access-Control-Allow-Credentials", "true")
	// ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Home Page",
	})
}

func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
	// origin := ctx.GetHeader("Origin")
	// if allowed := CheckOrigin(origin); allowed != "" {
	// 	ctx.Header("Access-Control-Allow-Origin", allowed)
	// } else {
	// 	ctx.AbortWithStatus(http.StatusForbidden)
	// 	return
	// }

	// ctx.Header("Access-Control-Allow-Credentials", "true")
	// ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
	// ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept")
	// log.Println("Received contact form request")
	// log.Printf("Contact form endpoint hit with path: %s", ctx.Request.URL.Path)
	// log.Printf("Request headers: %v", ctx.Request.Header)
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

	// log.Printf("Form data received: %+v", form)
	// log.Printf("reCAPTCHA token received: %s", form.Token)

	if !verifyRecaptcha(form.Token) {
		// log.Printf("reCAPTCHA verification failed for token: %s", form.Token)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
		return
	}

	err := email.SendEmail(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_TO"),
		form.Subject,
		"Name: "+form.Name+"\nEmail: "+form.Email+"\n\nMessage:\n"+form.Message,
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

// func sendEmail(name, email, subject, body, smtpPassword string) error {

// 	// log.Printf("SMTP Config - Host: %s, Port: %s, From: %s, To: %s",
// 	// 	os.Getenv("SMTP_HOST"),
// 	// 	os.Getenv("SMTP_PORT"),
// 	// 	os.Getenv("SMTP_FROM"),
// 	// 	os.Getenv("SMTP_TO"))

// 	// log.Printf("SMTP_PASSWORD length: %d", len(smtpPassword))

// 	m := gomail.NewMessage()
// 	m.SetHeader("From", os.Getenv("SMTP_FROM"))
// 	m.SetHeader("To", os.Getenv("SMTP_TO"))
// 	m.SetHeader("Reply-To", email)
// 	m.SetHeader("Subject", subject)
// 	m.SetBody("text/html", "Name: "+name+"\nEmail: "+email+"\n\nMessage:\n"+body)

// 	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
// 	if err != nil {
// 		log.Printf("Error converting SMTP_PORT to int: %v", err)
// 		return fmt.Errorf("invalid SMTP_PORT configuration")
// 	}

// 	d := gomail.NewDialer(
// 		os.Getenv("SMTP_HOST"),
// 		port,
// 		os.Getenv("SMTP_USERNAME"),
// 		smtpPassword,
// 	)
// 	d.TLSConfig = &tls.Config{
// 		ServerName: os.Getenv("SMTP_HOST"),
// 	}

// 	// log.Printf("Using SMTP credentials: %s / %s", os.Getenv("SMTP_USERNAME"), smtpPassword)

// 	// log.Println("Sending...")
// 	if err := d.DialAndSend(m); err != nil {
// 		log.Fatalf("Send failed: %v", err)
// 	}
// 	// log.Println("Sent.")
// 	return nil
// }

func verifyRecaptcha(token string) bool {
	secret := os.Getenv("RECAPTCHA_SECRET")
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	// log.Printf("Verifying reCAPTCHA with secret: %s, token: %s", secret, token)
	// log.Printf("Using RECAPTCHA_SECRET: %s", secret)

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

	// log.Printf("reCAPTCHA verification result: %+v", result)
	// log.Printf("Using RECAPTCHA_SECRET: %s", secret)
	return result.Success
}
