package tests

import (
	"testing"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/stretchr/testify/assert"
)

func TestUserModel(t *testing.T) {
	user := models.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		Timezone:     "UTC",
	}

	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.PasswordHash)
	assert.Equal(t, "user", user.Role)
	assert.Equal(t, "UTC", user.Timezone)
}

func TestBlogModel(t *testing.T) {
	blog := models.Blog{
		Title:       "Test Blog",
		FullContent: "This is test content",
		Excerpt:     "Test excerpt",
		Category:    "fitness",
		IsFeatured:  true,
	}

	assert.Equal(t, "Test Blog", blog.Title)
	assert.Equal(t, "This is test content", blog.FullContent)
	assert.Equal(t, "Test excerpt", blog.Excerpt)
	assert.Equal(t, "fitness", blog.Category)
	assert.True(t, blog.IsFeatured)
}

func TestWorkoutProgramModel(t *testing.T) {
	program := models.WorkoutProgram{
		Name:        "Test Workout Program",
		Description: "This is a test workout program",
		Duration:    30,
		Difficulty:  "beginner",
		IsActive:    true,
	}

	assert.Equal(t, "Test Workout Program", program.Name)
	assert.Equal(t, "This is a test workout program", program.Description)
	assert.Equal(t, 30, program.Duration)
	assert.Equal(t, "beginner", program.Difficulty)
	assert.True(t, program.IsActive)
}

func TestNewsletterSubscriberModel(t *testing.T) {
	subscriber := models.NewsletterSubscriber{
		Email:        "newsletter@example.com",
		IsConfirmed:  true,
		ConfirmToken: "test-token",
	}

	assert.Equal(t, "newsletter@example.com", subscriber.Email)
	assert.True(t, subscriber.IsConfirmed)
	assert.Equal(t, "test-token", subscriber.ConfirmToken)
}

func TestJobModel(t *testing.T) {
	job := models.Job{
		SessionID:     "session_123",
		CustomerEmail: "test@example.com",
		Status:        "pending",
		Attempts:      0,
		LastAttempt:   time.Now(),
	}

	assert.Equal(t, "session_123", job.SessionID)
	assert.Equal(t, "test@example.com", job.CustomerEmail)
	assert.Equal(t, "pending", job.Status)
	assert.Equal(t, 0, job.Attempts)
	assert.False(t, job.LastAttempt.IsZero())
}
