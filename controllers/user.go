package controllers

import (
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// func (uc *UserController) GetSignupForm(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Registration form"})
// }

// func (uc *UserController) SignupUser(ctx *gin.Context) {
// 	if ctx.Request.Method == "OPTIONS" {
// 		return
// 	}

// 	var payload error
// 	var user models.User

// 	if err := ctx.ShouldBindJSON(&user); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
// 		log.Println("Payload binding error:", err)
// 		return
// 	}

// 	log.Println("Received payload:", payload)

// 	var existingUser models.User
// 	result := uc.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		log.Println("User not found, proceeding to create.")
// 		// Proceed to create the user
// 	} else if result.Error != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
// 		log.Println("Database error:", result.Error)
// 		return
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	user.Password = string(hash)
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	if err := uc.DB.Create(&user).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"message": "User created successfully",
// 		"user": gin.H{
// 			"id":       user.ID,
// 			"username": user.Username,
// 			"email":    user.Email,
// 			"isAdmin":  user.IsAdmin,
// 		},
// 	})
// }

// func (uc *UserController) GetLoginForm(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Login form"})
// }

// func (uc *UserController) LoginUser(ctx *gin.Context) {
// 	log.Printf("LoginUser endpoint hit")
// 	var loginRequest struct {
// 		Email    string `json:"email" binding:"required"`
// 		Password string `json:"password" binding:"required"`
// 	}

// 	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
// 		return
// 	}

// 	var user models.User
// 	uc.DB.First(&user, "email = ?", loginRequest.Email)
// 	if user.ID == 0 {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email"})
// 		return
// 	}

// 	log.Println("Comparing passwords...")
// 	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
// 	if err != nil {
// 		log.Println("Password comparison failed:", err)
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
// 		return
// 	}

// 	accessToken, err := middleware.GenerateToken(user)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
// 		return
// 	}

// 	refreshToken, err := middleware.GenerateRefreshToken(user)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
// 		return
// 	}

// 	domain, secure, httpOnly, err := middleware.GetCookieSettings()
// 	if err != nil {
// 		log.Fatalf("Failed to parse environment variables: %v", err)
// 	}

// 	ctx.SetSameSite(http.SameSiteLaxMode)
// 	ctx.SetCookie(
// 		"access_token",
// 		accessToken,
// 		3600*1,
// 		"/",
// 		domain,
// 		secure,
// 		httpOnly,
// 	)

// 	ctx.SetSameSite(http.SameSiteLaxMode)
// 	ctx.SetCookie(
// 		"refresh_token",
// 		refreshToken,
// 		3600*24*30,
// 		"/",
// 		domain,
// 		secure,
// 		httpOnly,
// 	)

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Login successful",
// 		"user": gin.H{
// 			"_id":      user.ID,
// 			"email":    user.Email,
// 			"username": user.Username,
// 			"isAdmin":  user.IsAdmin,
// 		},
// 	})
// }

// func (uc *UserController) ForgotPassword(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Reset Password form"})
// }

// func (uc *UserController) ResetPassword(ctx *gin.Context) {
// 	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent"})
// }

// func (uc *UserController) LogoutUser(ctx *gin.Context) {
// 	log.Println("LogoutUser endpoint hit")

// 	domain, err := middleware.GetLogoutCookieSettings()
// 	if err != nil {
// 		log.Fatalf("Failed to parse environment variables: %v", err)
// 		log.Println("GetLogoutCookieSettings error:", err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process logout"})
// 		return
// 	}

// 	secure, httpOnly := false, false
// 	log.Printf("Domain: %s, Secure: %v, HttpOnly: %v", domain, secure, httpOnly)

// 	ctx.SetSameSite(http.SameSiteLaxMode)
// 	ctx.SetCookie(
// 		"access_token",
// 		"",
// 		-1,
// 		"/",
// 		domain,
// 		secure,
// 		httpOnly,
// 	)
// 	log.Println("Access token cookie cleared")

// 	ctx.SetSameSite(http.SameSiteLaxMode)
// 	ctx.SetCookie(
// 		"refresh_token",
// 		"",
// 		-1,
// 		"/",
// 		domain,
// 		secure,
// 		httpOnly,
// 	)
// 	log.Println("Refresh token cookie cleared")

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
// }

// func (uc *UserController) GetProfile(ctx *gin.Context) {
// 	if ctx.Request.Method == "OPTIONS" {
// 		return
// 	}

// 	userID := ctx.Param("id")
// 	if userID == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
// 		return
// 	}

// 	var user models.User
// 	err := uc.DB.Where("id = ?", userID).First(&user).Error
// 	if err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"_id":      user.ID,
// 		"username": user.Username,
// 		"email":    user.Email,
// 		"isAdmin":  user.IsAdmin,
// 	})
// }
