package controllers

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type HomeController struct {
// 	DB *gorm.DB
// }

// func NewHomeController(db *gorm.DB) *HomeController {
// 	return &HomeController{DB: db}
// }

// type ContactForm struct {
// 	Name    string `json:"name" binding:"required"`
// 	Email   string `json:"email" binding:"required,email"`
// 	Message string `json:"message" binding:"required"`
// 	Token   string `json:"token"`
// }

// func (hc *HomeController) GetHome(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Welcome to the Home Page",
// 	})
// }

// // func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
// // 	var form ContactForm
// // 	if err := ctx.ShouldBindJSON(&form); err != nil {
// // 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// // 		return
// // 	}

// // 	if !verifyRecaptcha(form.Token) {
// // 		ctx.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
// // 		return
// // 	}

// // 	err := sendEmail(form.Name, form.Email, form.Message)
// // 	if err != nil {
// // 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
// // 		return
// // 	}

// // 	// Here you can process the form submission, e.g., send an email or store in a database
// // 	log.Printf("Contact form submitted: %+v\n", form)

// // 	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message received!"})
// // }

// // func sendEmail(name, email, message string) error {
// // 	m := gomail.NewMessage()
// // 	m.SetHeader("From", os.Getenv("SMTP_FROM")) // Sender email from environment variable
// // 	m.SetHeader("To", os.Getenv("SMTP_TO"))     // Recipient email
// // 	m.SetHeader("Subject", "New Contact Form Submission")
// // 	m.SetBody("text/plain", "Name: "+name+"\nEmail: "+email+"\n\nMessage:\n"+message)

// // 	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

// // 	d := gomail.NewDialer(
// // 		os.Getenv("SMTP_HOST"),
// // 		port,
// // 		os.Getenv("SMTP_USERNAME"),
// // 		os.Getenv("SMTP_PASSWORD"),
// // 	)

// // 	return d.DialAndSend(m)
// // }

// // func verifyRecaptcha(token string) bool {
// // 	secret := os.Getenv("RECAPTCHA_SECRET")
// // 	url := "https://www.google.com/recaptcha/api/siteverify"

// // 	data := map[string]string{
// // 		"secret":   secret,
// // 		"response": token,
// // 	}
// // 	jsonData, _ := json.Marshal(data)

// // 	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
// // 	if err != nil {
// // 		return false
// // 	}
// // 	defer resp.Body.Close()

// // 	var recaptchaResp struct {
// // 		Success bool `json:"success"`
// // 	}
// // 	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
// // 		return false
// // 	}

// // 	return recaptchaResp.Success
// // }
