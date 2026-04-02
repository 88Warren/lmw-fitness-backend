package workers

import (
	"log"
	"os"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"github.com/88warren/lmw-fitness-backend/utils/email"
	"github.com/88warren/lmw-fitness-backend/utils/emailtemplates"
	"gorm.io/gorm"
)

// StartReminderWorker runs once a day and emails users who haven't worked out in 2+ days
func StartReminderWorker(db *gorm.DB) {
	log.Println("Reminder worker started")

	// Run immediately on startup (catches any missed sends), then every 24h
	go func() {
		sendReminders(db)
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			sendReminders(db)
		}
	}()
}

func sendReminders(db *gorm.DB) {
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		// Try Kubernetes secret path
		if data, err := os.ReadFile("/etc/secrets/smtp-password"); err == nil {
			smtpPassword = string(data)
		}
	}
	if smtpPassword == "" {
		log.Println("Reminder worker: SMTP_PASSWORD not set, skipping")
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://lmwfitness.co.uk"
	}
	fromAddress := os.Getenv("SMTP_FROM")

	// Find users who:
	// - have worked out at least once (LastWorkoutDate is not null)
	// - haven't worked out in 2+ days
	// - haven't opted out of reminders
	twoDaysAgo := time.Now().Add(-48 * time.Hour)

	var users []models.User
	if err := db.Where(
		"last_workout_date IS NOT NULL AND last_workout_date < ? AND reminder_opt_out = false",
		twoDaysAgo,
	).Find(&users).Error; err != nil {
		log.Printf("Reminder worker: failed to query users: %v", err)
		return
	}

	log.Printf("Reminder worker: found %d users to remind", len(users))

	for _, u := range users {
		daysSince := int(time.Since(*u.LastWorkoutDate).Hours() / 24)

		// Only send on day 2, day 5, and day 10 to avoid spamming
		if daysSince != 2 && daysSince != 5 && daysSince != 10 {
			continue
		}

		body := emailtemplates.GenerateWorkoutReminderEmailBody(
			u.Email,
			daysSince,
			u.CurrentStreak,
			frontendURL,
		)

		err := email.SendEmail(
			fromAddress,
			u.Email,
			"Your workout is waiting 💪",
			body,
			"",
			smtpPassword,
		)
		if err != nil {
			log.Printf("Reminder worker: failed to send to %s: %v", u.Email, err)
		} else {
			log.Printf("Reminder worker: sent reminder to %s (%d days inactive)", u.Email, daysSince)
		}
	}
}
