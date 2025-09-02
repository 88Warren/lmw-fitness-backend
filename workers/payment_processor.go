package workers

import (
	"log"
	"time"

	"github.com/88warren/lmw-fitness-backend/models"
	"gorm.io/gorm"
)

type PaymentProcessor interface {
	ProcessPaymentSuccess(sessionID string, customerEmail string) error
}

type JobProcessor struct {
	db       *gorm.DB
	pc       PaymentProcessor
	jobChan  chan struct{}
	stopChan chan struct{}
}

func NewJobProcessor(db *gorm.DB, pc PaymentProcessor) *JobProcessor {
	return &JobProcessor{
		db:       db,
		pc:       pc,
		jobChan:  make(chan struct{}, 100),
		stopChan: make(chan struct{}),
	}
}

func StartPaymentWorker(db *gorm.DB, pc PaymentProcessor) {
	processor := NewJobProcessor(db, pc)

	SetGlobalProcessor(processor)

	// log.Printf("=== PAYMENT WORKER STARTED ===")
	// log.Printf("Worker is ready to process jobs immediately")

	processor.processPendingJobs()

	go processor.workerLoop()

	select {}
}

func (jp *JobProcessor) workerLoop() {
	fallbackTicker := time.NewTicker(5 * time.Minute)
	defer fallbackTicker.Stop()

	for {
		select {
		case <-jp.jobChan:
			// log.Printf("=== Processing jobs triggered by event ===")
			jp.processPendingJobs()

		case <-fallbackTicker.C:
			// log.Printf("=== Fallback check for missed jobs ===")
			jp.processPendingJobs()

		case <-jp.stopChan:
			// log.Printf("=== PAYMENT WORKER STOPPED ===")
			return
		}
	}
}

func (jp *JobProcessor) TriggerJobProcessing() {
	select {
	case jp.jobChan <- struct{}{}:
		// log.Printf("Job processing triggered")
	default:
		// log.Printf("Job processing channel full, will process on next fallback cycle")
	}
}

func (jp *JobProcessor) Stop() {
	close(jp.stopChan)
}

func (jp *JobProcessor) processPendingJobs() {
	var jobs []models.Job
	result := jp.db.Where("status = ? AND attempts < ?", "pending", 5).Find(&jobs)

	if result.Error != nil {
		log.Printf("Error querying for pending jobs: %v", result.Error)
		return
	}

	// var allJobs []models.Job
	// if allResult := jp.db.Find(&allJobs); allResult.Error == nil {
	// 	log.Printf("Total jobs in database: %d", len(allJobs))
	// 	for _, job := range allJobs {
	// 		log.Printf("Job: ID=%d, Session=%s, Email=%s, Status=%s, Attempts=%d",
	// 			job.ID, job.SessionID, job.CustomerEmail, job.Status, job.Attempts)
	// 	}
	// }

	if len(jobs) == 0 {
		// log.Printf("No pending jobs found")
		return
	}

	// log.Printf("Found %d pending jobs", len(jobs))

	for _, job := range jobs {
		// log.Printf("=== Processing job for session: %s ===", job.SessionID)
		// log.Printf("Job details: ID=%d, Email=%s, Status=%s, Attempts=%d",
		// 	job.ID, job.CustomerEmail, job.Status, job.Attempts)

		job.Status = "processing"
		if err := jp.db.Save(&job).Error; err != nil {
			log.Printf("Error updating job status to processing: %v", err)
			continue
		}

		err := jp.pc.ProcessPaymentSuccess(job.SessionID, job.CustomerEmail)
		if err != nil {
			log.Printf("Job for session %s failed: %v", job.SessionID, err)
			job.Status = "failed"
			job.Attempts++
			if err := jp.db.Save(&job).Error; err != nil {
				log.Printf("Error updating failed job: %v", err)
			}
		} else {
			job.Status = "completed"
			if err := jp.db.Save(&job).Error; err != nil {
				log.Printf("Error updating completed job: %v", err)
			} else {
				// log.Printf("Job for session %s completed successfully.", job.SessionID)
			}
		}
	}
}

var globalProcessor *JobProcessor

func GetGlobalProcessor() *JobProcessor {
	return globalProcessor
}

func SetGlobalProcessor(processor *JobProcessor) {
	globalProcessor = processor
}
