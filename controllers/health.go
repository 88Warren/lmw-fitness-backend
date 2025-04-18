package controllers

import (
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
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "Ready"})
}

func (healthController *HealthController) checkDatabaseConnection() error {
	sqlDB, err := healthController.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
