package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthController struct {
	DB *gorm.DB
}

func NewHealthController(db *gorm.DB) *HealthController {
	return &HealthController{DB: db}
}

func (healthController *HealthController) LivenessCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "Live"})
}

func (healthController *HealthController) ReadinessCheck(ctx *gin.Context) {
	err := healthController.checkDatabaseConnection()
	if err != nil {
		log.Printf("Readiness check failed: %v", err)
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "Not Ready",
			"error":  fmt.Sprintf("Database connection failed: %v", err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "Ready"})
}

func (healthController *HealthController) checkDatabaseConnection() error {
	if healthController.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := healthController.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}

	return nil
}
