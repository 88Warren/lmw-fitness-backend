package workers

import (
	"log"
	"time"

	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

func StartPaymentWorker(db *gorm.DB, pc *controllers.PaymentController) {
	log.Printf("=== PAYMENT WORKER STARTED ===")
	log.Printf("Worker will check for pending jobs every 5 seconds")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var jobs []models.Job
		result := db.Where("status = ? AND attempts < ?", "pending", 5).Find(&jobs)

		if result.Error != nil {
			log.Printf("Error querying for pending jobs: %v", result.Error)
			continue
		}

		log.Printf("Found %d pending jobs", len(jobs))

		for _, job := range jobs {
			log.Printf("=== Processing job for session: %s ===", job.SessionID)
			log.Printf("Job details: ID=%d, Email=%s, Status=%s, Attempts=%d",
				job.ID, job.CustomerEmail, job.Status, job.Attempts)

			job.Status = "processing"
			if err := db.Save(&job).Error; err != nil {
				log.Printf("Error updating job status to processing: %v", err)
				continue
			}

			err := pc.ProcessPaymentSuccess(job.SessionID, job.CustomerEmail)
			if err != nil {
				log.Printf("Job for session %s failed: %v", job.SessionID, err)
				job.Status = "failed"
				job.Attempts++
				if err := db.Save(&job).Error; err != nil {
					log.Printf("Error updating failed job: %v", err)
				}
			} else {
				job.Status = "completed"
				if err := db.Save(&job).Error; err != nil {
					log.Printf("Error updating completed job: %v", err)
				} else {
					log.Printf("Job for session %s completed successfully.", job.SessionID)
				}
			}
		}

		if len(jobs) == 0 {
			log.Printf("No pending jobs found, waiting...")
		}
	}
}
