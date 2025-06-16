package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsletterController struct {
	DB *gorm.DB
}

func NewNewsletterController(db *gorm.DB) *NewsletterController {
	return &NewsletterController{DB: db}
}

type BrevoAddContactResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
}

type BrevoErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SubscribeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (nc *NewsletterController) Subscribe(ctx *gin.Context) {
	var req SubscribeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	brevoAPIKey := os.Getenv("BREVO_API_KEY")
	brevoListID := os.Getenv("BREVO_LIST_ID")
	brevoAPIURL := os.Getenv("BREVO_API_URL")
	brevoDOIRedirectURL := os.Getenv("BREVO_DOI_REDIRECT_URL")
	brevoDOITemplateID := os.Getenv("BREVO_DOI_TEMPLATE_ID")

	if brevoAPIKey == "" || brevoListID == "" || brevoAPIURL == "" {
		log.Println("Brevo API environment variables not set. Cannot subscribe.")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Newsletter service not configured."})
		return
	}

	url := fmt.Sprintf("%s/contacts", brevoAPIURL)

	requestBody := map[string]interface{}{
		"email": req.Email,
		"attributes": map[string]string{
			"SMS":       "",
			"FIRSTNAME": "",
			"LASTNAME":  "",
		},
		"listIds":          []int{atoi(brevoListID)},
		"updateEnabled":    true,
		"emailBlacklisted": false,
		"smsBlacklisted":   false,
		"status":           "pending",
		"templateId":       atoi(brevoDOITemplateID),
		"redirectionUrl":   brevoDOIRedirectURL,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshaling Brevo request body: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error."})
		return
	}

	log.Printf("Brevo Request URL: %s", url)
	log.Printf("Brevo Request Body: %s", string(jsonBody))
	log.Printf("Using API Key (first 5 chars): %s...", brevoAPIKey[:5])

	client := &http.Client{Timeout: 10 * time.Second}
	reqAPI, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating Brevo API request: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe."})
		return
	}

	reqAPI.Header.Set("Content-Type", "application/json")
	reqAPI.Header.Set("api-key", brevoAPIKey)

	resp, err := client.Do(reqAPI)
	if err != nil {
		log.Printf("Error sending request to Brevo API: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe. Network error."})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Subscription successful! Please check your inbox to confirm."})
	} else {
		var errorBody BrevoErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorBody); err != nil {
			log.Printf("Error decoding Brevo error response: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe. Unknown error from service."})
			return
		}
		log.Printf("Brevo API error (Status: %d, Code: %s): %s", resp.StatusCode, errorBody.Code, errorBody.Message)

		if errorBody.Code == "duplicate_parameter" || errorBody.Code == "already_exist" {
			ctx.JSON(http.StatusOK, gin.H{"message": "You are already subscribed to our newsletter! Please check your inbox for confirmation."})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to subscribe: %s", errorBody.Message)})
		}
	}
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Error converting string to int: %v", err)
		return 0
	}
	return i
}
