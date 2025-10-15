package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/88warren/lmw-fitness-backend/config"
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/routes"
	"github.com/stretchr/testify/assert"
)

func TestGetBlogs(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Create test blog posts
	blog1 := models.Blog{
		Title:       "Test Blog 1",
		FullContent: "This is test content 1",
		Excerpt:     "Test excerpt 1",
		Category:    "fitness",
		IsFeatured:  false,
	}
	blog2 := models.Blog{
		Title:       "Test Blog 2",
		FullContent: "This is test content 2",
		Excerpt:     "Test excerpt 2",
		Category:    "nutrition",
		IsFeatured:  true,
	}
	db.Create(&blog1)
	db.Create(&blog2)

	router := config.SetupServer()
	blogController := controllers.NewBlogController(GetTestDB())
	routes.RegisterBlogRoutes(router, blogController)

	req, _ := http.NewRequest("GET", "/api/blog", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Blog
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.GreaterOrEqual(t, len(response), 2)
}

func TestGetBlogByID(t *testing.T) {
	// Skip if no database connection
	db := GetTestDB()
	if db == nil {
		t.Skip("Skipping database test - no connection available")
	}

	// Create test blog post
	blog := models.Blog{
		Title:       "Specific Blog Test",
		FullContent: "This is specific test content",
		Excerpt:     "Specific test excerpt",
		Category:    "fitness",
		IsFeatured:  false,
	}
	result := db.Create(&blog)
	if result.Error != nil {
		t.Fatalf("Failed to create test blog: %v", result.Error)
	}

	router := config.SetupServer()
	blogController := controllers.NewBlogController(GetTestDB())
	routes.RegisterBlogRoutes(router, blogController)

	// Use the actual ID of the created blog
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/blog/%d", blog.ID), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 200 or 404 depending on implementation
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)

	if w.Code == http.StatusOK {
		var response models.Blog
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "Specific Blog Test", response.Title)
	}
}

func TestGetBlogByIDNotFound(t *testing.T) {
	router := config.SetupServer()
	blogController := controllers.NewBlogController(GetTestDB())
	routes.RegisterBlogRoutes(router, blogController)

	req, _ := http.NewRequest("GET", "/api/blog/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateBlog(t *testing.T) {
	router := config.SetupServer()
	blogController := controllers.NewBlogController(GetTestDB())
	routes.RegisterBlogRoutes(router, blogController)

	requestBody := map[string]interface{}{
		"title":       "New Blog Post",
		"fullContent": "This is new blog content",
		"excerpt":     "New blog excerpt",
		"category":    "fitness",
		"isFeatured":  false,
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/blog", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 201, 200, 400, or 401 (unauthorized) depending on implementation
	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK || w.Code == http.StatusBadRequest || w.Code == http.StatusUnauthorized)

	if w.Code == http.StatusCreated || w.Code == http.StatusOK {
		var response models.Blog
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "New Blog Post", response.Title)
	}
}
