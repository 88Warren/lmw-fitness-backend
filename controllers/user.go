package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/utils/email"
	"github.com/88warren/lmw-fitness-backend/utils/emailtemplates"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

var (
	hasUpperCase   = regexp.MustCompile(`[A-Z]`)
	hasSpecialChar = regexp.MustCompile(`[!@#$^&*]`)
)

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ValidatePassword(req.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user with this email already exists
	var existingUser models.User
	if err := uc.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	} else if err != gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking user existence"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	if result := uc.DB.Create(&user); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user: " + result.Error.Error()})
		return
	}

	// Generate JWT token for the newly registered user
	token, err := GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if result := uc.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during login"})
		return
	}

	log.Printf("Stored Hashed Password for %s: %s", user.Email, user.PasswordHash)
	log.Printf("Login attempt plaintext password: %s", req.Password)

	// Compare provided password with hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": models.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	// Find user by ID
	if result := uc.DB.First(&user, userID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile: " + result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

func GenerateJWT(userID uint, email, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// log.Println("JWT_SECRET environment variable not set. Using a default (NOT SECURE FOR PRODUCTION!).")
		jwtSecret = "supersecretjwtkey"
	}

	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"email":      email,
		"role":       role,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uc *UserController) RequestPasswordReset(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := uc.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during password reset request"})
		return
	}

	// 1. Generate a unique, time-limited token
	token, err := generateSecureToken(32)
	if err != nil {
		log.Printf("Error generating reset token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token."})
		return
	}

	expiresAt := time.Now().Add(time.Hour * 12)

	// 2. Save token to database, invalidate any existing tokens for this user
	uc.DB.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{})

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	if result := uc.DB.Create(&resetToken); result.Error != nil {
		log.Printf("Error saving reset token to DB: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token."})
		return
	}

	// 3. Send email with reset link
	resetLink := fmt.Sprintf("%s/reset-password/%s", os.Getenv("ALLOWED_ORIGIN"), token)
	emailSubject := "LMW Fitness - Password Reset Request"
	emailBody := emailtemplates.GeneratePasswordResetEmailBody(user.Email, resetLink)

	smtpPassword := getSMTPPasswordFromSecrets()

	if err := email.SendEmail(
		os.Getenv("SMTP_FROM"),
		user.Email,
		emailSubject,
		emailBody,
		"",
		smtpPassword,
	); err != nil {
		log.Printf("Error sending password reset email: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "If an account with that email exists, a password reset link has been sent."})
}

func (uc *UserController) VerifyResetToken(ctx *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resetToken models.PasswordResetToken
	if result := uc.DB.Where("token = ?", req.Token).First(&resetToken); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during token verification."})
		return
	}

	if time.Now().After(resetToken.ExpiresAt) {
		uc.DB.Delete(&resetToken)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password reset token has expired. Please request a new one."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Token is valid."})
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if !hasUpperCase.MatchString(password) {
		return fmt.Errorf("password must contain at least one capital letter")
	}

	if !hasSpecialChar.MatchString(password) {
		return fmt.Errorf("password must contain at least one special character (!@#$^&*)")
	}

	return nil
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Password complexity validation (Perform this early)
	if err := ValidatePassword(req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find and validate the token
	var resetToken models.PasswordResetToken
	if result := uc.DB.Where("token = ?", req.Token).First(&resetToken); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired password reset link."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error during password reset."})
		return
	}

	if time.Now().After(resetToken.ExpiresAt) {
		uc.DB.Delete(&resetToken)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password reset link has expired. Please request a new one."})
		return
	}

	// 3. Find the user associated with the token
	var user models.User
	if result := uc.DB.First(&user, resetToken.UserID); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User associated with this token not found."})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding user for password reset."})
		return
	}

	// 4. Prevent using the same password as the old one
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.NewPassword)); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "New password cannot be the same as your old password."})
		return
	} else if err != bcrypt.ErrMismatchedHashAndPassword {
		log.Printf("Bcrypt comparison error during password reset: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred during password validation."})
		return
	}

	// 5. Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password."})
		return
	}

	log.Printf("New Hashed Password for user %d: %s", user.ID, string(hashedPassword))

	// 6. Update the user's password in the database
	user.PasswordHash = string(hashedPassword)
	if result := uc.DB.Save(&user); result.Error != nil {
		log.Printf("Error updating user password: %v", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password."})
		return
	}

	// 7. Invalidate (delete) the token after successful use
	uc.DB.Delete(&resetToken)

	ctx.JSON(http.StatusOK, gin.H{"message": "Your password has been reset successfully!"})
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func getSMTPPasswordFromSecrets() string {
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		return os.Getenv("SMTP_PASSWORD")
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

	return string(secret.Data["SMTP_PASSWORD"])
}
