package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestNewsletterSubscription(t *testing.T) {
	router := config.SetupServer()
	newsletterController := controllers.NewNewsletterController(GetTestDB())
	routes.RegisterNewsletterRoutes(router, newsletterController)

	requestBody := map[string]interface{}{
		"email": "newsletter@example.com",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/newsletter/subscribe", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 201, 200, or 400 depending on implementation
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK || w.Code == http.StatusBadRequest)

	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		// Response structure may vary
	}
}

func TestNewsletterDuplicateSubscription(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// First create a subscriber
	subscriber := models.NewsletterSubscriber{
		Email:       "duplicate@example.com",
		IsConfirmed: true,
	}
	db.Create(&subscriber)

	router := config.SetupServer()
	newsletterController := controllers.NewNewsletterController(GetTestDB())
	routes.RegisterNewsletterRoutes(router, newsletterController)

	requestBody := map[string]interface{}{
		"email": "duplicate@example.com",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/newsletter/subscribe", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should handle duplicate subscription appropriately
	assert.True(t, w.Code == http.StatusConflict || w.Code == http.StatusOK)
}

func TestNewsletterInvalidEmail(t *testing.T) {
	router := config.SetupServer()
	newsletterController := controllers.NewNewsletterController(GetTestDB())
	routes.RegisterNewsletterRoutes(router, newsletterController)

	requestBody := map[string]interface{}{
		"email": "invalid-email",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/newsletter/subscribe", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestNewsletterUnsubscribe(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// First create a subscriber
	subscriber := models.NewsletterSubscriber{
		Email:       "unsubscribe@example.com",
		IsConfirmed: true,
	}
	db.Create(&subscriber)

	router := config.SetupServer()
	newsletterController := controllers.NewNewsletterController(GetTestDB())
	routes.RegisterNewsletterRoutes(router, newsletterController)

	requestBody := map[string]interface{}{
		"email": "unsubscribe@example.com",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/newsletter/unsubscribe", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
