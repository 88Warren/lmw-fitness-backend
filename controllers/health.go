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

func (hc *HealthController) LivenessCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Liveness check passed"})
}

func (hc *HealthController) ReadinessCheck(ctx *gin.Context) {
	err := hc.checkDatabaseConnection()
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Readiness check passed"})
}

func (hc *HealthController) checkDatabaseConnection() error {
	sqlDB, err := hc.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
