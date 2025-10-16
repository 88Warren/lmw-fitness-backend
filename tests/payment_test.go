package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestCreateCheckoutSession(t *testing.T) {
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		t.Skip("Skipping Stripe test - no API key")
	}

	router := config.SetupServer()
	paymentController := controllers.NewPaymentController(GetTestDB())
	routes.RegisterPaymentRoutes(router, paymentController)

	requestBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"priceId":  "price_test_123",
				"quantity": 1,
			},
		},
		"customerEmail": "info@lmwfitness.co.uk",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/create-checkout-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// This should still return 400 because we're using a test Stripe key
	// but now it won't fail on the binding validation
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPaymentValidation(t *testing.T) {
	if os.Getenv("STRIPE_SECRET_KEY") == "" {
		t.Skip("Skipping Stripe test - no API key")
	}

	router := config.SetupServer()
	paymentController := controllers.NewPaymentController(GetTestDB())
	routes.RegisterPaymentRoutes(router, paymentController)

	// Test with missing items (should fail validation)
	requestBody := map[string]interface{}{
		"customerEmail": "test@example.com",
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/create-checkout-session", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Items")
}
