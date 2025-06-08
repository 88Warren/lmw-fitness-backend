package middleware

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	// allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	// log.Println("allowedOrigin: ", allowedOrigin)
	// allowCredentialsStr := os.Getenv("ALLOW_CREDENTIALS")
	// allowCredentials, err := strconv.ParseBool(allowCredentialsStr)
	// if err != nil {
	// log.Printf("Warning: ALLOW_CREDENTIALS environment variable '%s' is not a valid boolean. Defaulting to false.", allowCredentialsStr)
	// 	allowCredentials = false
	// }
	// log.Println("AllowCredentials (parsed from .env): ", allowCredentials)

	// return cors.New(cors.Config{
	// 	AllowOriginFunc: func(origin string) bool {
	// 		if allowCredentials {
	// 			return origin == allowedOrigin
	// 		}
	// 		return origin == allowedOrigin || allowedOrigin == "*"
	// 	},

	envAllowedOrigins := os.Getenv("ALLOWED_ORIGIN")
	var allowedOriginsList []string
	if envAllowedOrigins == "" {
		// log.Println("WARNING: ALLOWED_ORIGIN environment variable is not set for CORS. No origins will be allowed by default.")
		allowedOriginsList = []string{}
	} else {
		allowedOriginsList = strings.Split(envAllowedOrigins, ",")
		// Trim spaces from each origin
		for i, origin := range allowedOriginsList {
			allowedOriginsList[i] = strings.TrimSpace(origin)
		}
		// log.Printf("CORS middleware configured with allowed origins: %v", allowedOriginsList)
	}

	allowCredentialsStr := os.Getenv("ALLOW_CREDENTIALS")
	allowCredentials, err := strconv.ParseBool(allowCredentialsStr)
	if err != nil {
		// log.Printf("Warning: ALLOW_CREDENTIALS environment variable '%s' is not a valid boolean. Defaulting to false.", allowCredentialsStr)
		allowCredentials = false
	}
	// log.Println("AllowCredentials (parsed from .env): ", allowCredentials)

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range allowedOriginsList {
				if origin == allowed || allowed == "*" {
					return true
				}
			}
			return false
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Content-Disposition",
		},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
