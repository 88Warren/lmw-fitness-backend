package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Ensure you use v5
)

// AuthMiddleware is a Gin middleware to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		// Check if the token starts with "Bearer "
		if !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format. Must be 'Bearer <token>'"})
			ctx.Abort()
			return
		}

		// Extract the actual token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Get JWT secret from environment variables
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Println("JWT_SECRET environment variable not set. Using a default (NOT SECURE FOR PRODUCTION!).")
			jwtSecret = "supersecretjwtkey" // Fallback, should match the one in controllers/user_controller.go
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err) // Log the specific error
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		// Extract user ID, email, and role from claims
		userID, ok := claims["user_id"].(float64) // JWT numbers are float64 by default
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			ctx.Abort()
			return
		}
		email, ok := claims["email"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
			ctx.Abort()
			return
		}
		role, ok := claims["role"].(string)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
			ctx.Abort()
			return
		}

		// Set user information in the Gin context for subsequent handlers
		ctx.Set("userID", uint(userID)) // Convert back to uint
		ctx.Set("userEmail", email)
		ctx.Set("userRole", role)

		ctx.Next() // Proceed to the next handler in the chain
	}
}

// RoleMiddleware checks if the authenticated user has the required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			// This indicates AuthMiddleware didn't run or didn't set the role
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context. AuthMiddleware might be missing."})
			ctx.Abort()
			return
		}

		if userRole != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role permissions"})
			ctx.Abort()
			return
		}

		ctx.Next() // User has the required role, proceed
	}
}
