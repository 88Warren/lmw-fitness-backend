package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/laurawarren88/LMW_Fitness/models" // Import your models package
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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

// GetProfile fetches the profile of the authenticated user
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

	// Return user data, excluding sensitive fields like PasswordHash
	ctx.JSON(http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
}

func GenerateJWT(userID uint, email, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET environment variable not set. Using a default (NOT SECURE FOR PRODUCTION!).")
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
