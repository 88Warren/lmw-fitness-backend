package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	// Setup router
	router := config.SetupServer()

	// Setup health controller and routes
	healthController := controllers.NewHealthController(GetTestDB())
	routes.RegisterHealthRoutes(router, healthController)

	// Create test request
	req, _ := http.NewRequest("GET", "/api/health", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "status")
}

func TestDatabaseHealthCheck(t *testing.T) {
	db := GetTestDB()
	assert.NotNil(t, db)

	// Test database connection
	sqlDB, err := db.DB()
	assert.NoError(t, err)

	err = sqlDB.Ping()
	assert.NoError(t, err)
}
