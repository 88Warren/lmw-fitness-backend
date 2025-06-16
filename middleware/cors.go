package middleware

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	envAllowedOrigins := os.Getenv("ALLOWED_ORIGIN")
	var allowedOriginsList []string
	if envAllowedOrigins == "" {
		allowedOriginsList = []string{}
	} else {
		allowedOriginsList = strings.Split(envAllowedOrigins, ",")
		for i, origin := range allowedOriginsList {
			allowedOriginsList[i] = strings.TrimSpace(origin)
		}
		// log.Printf("CORS middleware configured with allowed origins: %v", allowedOriginsList)
	}

	allowCredentialsStr := os.Getenv("ALLOW_CREDENTIALS")
	allowCredentials, err := strconv.ParseBool(allowCredentialsStr)
	if err != nil {
		allowCredentials = false
	}

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// log.Printf("CORS Check: Incoming Origin = '%s'", origin)
			parsedOriginURL, err := url.Parse(origin)
			if err != nil {
				fmt.Printf("Error parsing origin URL '%s': %v\n", origin, err)
				return false
			}
			originScheme := parsedOriginURL.Scheme
			originHost := parsedOriginURL.Hostname()

			// log.Printf("CORS Check: Parsed Origin Scheme: '%s', Hostname: '%s'", originScheme, originHost)

			for _, allowed := range allowedOriginsList {
				// log.Printf("CORS Check: Comparing against Allowed Origin = '%s'", allowed)

				if strings.TrimSuffix(origin, "/") == strings.TrimSuffix(allowed, "/") || allowed == "*" {
					// log.Printf("CORS Check: Matched by exact string or global wildcard for '%s'", allowed)
					return true
				}

				parsedAllowedURL, err := url.Parse(allowed)
				if err != nil {
					log.Printf("CORS Check: Error parsing allowed origin URL '%s': %v", allowed, err)
					continue
				}

				if originScheme != parsedAllowedURL.Scheme {
					// log.Printf("CORS Check: Scheme mismatch. Origin: '%s', Allowed: '%s'", originScheme, parsedAllowedURL.Scheme)
					continue
				}
				allowedHostPattern := parsedAllowedURL.Hostname()

				if strings.HasPrefix(allowedHostPattern, "*.") {
					suffix := allowedHostPattern[1:]
					// log.Printf("CORS Check: Attempting wildcard match. Suffix = '%s'. Is '%s' suffix of '%s'? Result: %v", suffix, suffix, originHost, strings.HasSuffix(originHost, suffix))
					if strings.HasSuffix(originHost, suffix) {
						// log.Printf("CORS Check: Matched by wildcard domain: '%s'", allowed)
						return true
					}
				}
			}
			// fmt.Printf("Origin %s not allowed by CORS config\n", origin)
			return false
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type", "Content-Length", "Accept-Encoding", "Authorization",
			"Accept", "Origin", "Cache-Control", "X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length", "Content-Type", "Content-Disposition",
		},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	})
}
