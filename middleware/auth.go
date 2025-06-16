package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}
		if !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format. Must be 'Bearer <token>'"})
			ctx.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			// log.Println("JWT_SECRET environment variable not set. Using a default (NOT SECURE FOR PRODUCTION!).")
			jwtSecret = "supersecretjwtkey"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("Token validation error: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
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

		ctx.Set("userID", uint(userID))
		ctx.Set("userEmail", email)
		ctx.Set("userRole", role)

		ctx.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("userRole")
		if !exists {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found in context. AuthMiddleware might be missing."})
			ctx.Abort()
			return
		}

		if userRole != requiredRole {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role permissions"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
