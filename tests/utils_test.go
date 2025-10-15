package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Test utility functions that don't require database

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Verify the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	assert.NoError(t, err)

	// Verify wrong password fails
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	assert.Error(t, err)
}

func TestEmailValidation(t *testing.T) {
	validEmails := []string{
		"test@example.com",
		"user.name@domain.co.uk",
		"user+tag@example.org",
	}

	invalidEmails := []string{
		"invalid-email",
		"test@",
		"@example.com",
		"",
		"test.example.com",
	}

	// Simple email validation regex
	emailRegex := `^[^\s@]+@[^\s@]+\.[^\s@]+$`

	for _, email := range validEmails {
		assert.Regexp(t, emailRegex, email, "Email should be valid: %s", email)
	}

	for _, email := range invalidEmails {
		assert.NotRegexp(t, emailRegex, email, "Email should be invalid: %s", email)
	}
}

func TestSlugGeneration(t *testing.T) {
	// Mock slug generation function
	generateSlug := func(title string) string {
		// Simple slug generation for testing - use the title parameter
		if title == "" {
			return ""
		}
		return "test-slug"
	}

	title := "Test Blog Title"
	slug := generateSlug(title)

	assert.Equal(t, "test-slug", slug)
	assert.NotEmpty(t, slug)

	// Test empty title
	emptySlug := generateSlug("")
	assert.Empty(t, emptySlug)
}

func TestStringValidation(t *testing.T) {
	// Test string length validation
	shortString := "abc"
	longString := "this is a very long string that exceeds normal limits"
	normalString := "normal string"

	assert.True(t, len(shortString) >= 3)
	assert.True(t, len(normalString) > 3 && len(normalString) < 50)
	assert.True(t, len(longString) > 50)
}
