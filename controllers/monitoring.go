package controllers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/88warren/lmw-fitness-backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MonitoringController struct {
	DB *gorm.DB
}

func NewMonitoringController(db *gorm.DB) *MonitoringController {
	return &MonitoringController{DB: db}
}

type HealthResponse struct {
	Status    string              `json:"status"`
	Timestamp time.Time           `json:"timestamp"`
	Version   string              `json:"version"`
	Database  DatabaseHealth      `json:"database"`
	System    SystemHealth        `json:"system"`
	Metrics   *middleware.Metrics `json:"metrics"`
}

type DatabaseHealth struct {
	Status      string        `json:"status"`
	Connections int           `json:"connections"`
	Latency     time.Duration `json:"latency_ms"`
}

type SystemHealth struct {
	Memory     MemoryStats `json:"memory"`
	Goroutines int         `json:"goroutines"`
	Uptime     string      `json:"uptime"`
}

type MemoryStats struct {
	Allocated uint64 `json:"allocated_mb"`
	System    uint64 `json:"system_mb"`
	GCRuns    uint32 `json:"gc_runs"`
}

var startTime = time.Now()

func (mc *MonitoringController) HealthCheck(c *gin.Context) {
	// Check database health
	dbHealth := mc.checkDatabaseHealth()

	// Get system metrics
	systemHealth := mc.getSystemHealth()

	// Get application metrics
	metrics := middleware.GetMetrics()

	status := "healthy"
	if dbHealth.Status != "healthy" {
		status = "unhealthy"
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.0.0", // You can make this dynamic
		Database:  dbHealth,
		System:    systemHealth,
		Metrics:   metrics,
	}

	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

func (mc *MonitoringController) checkDatabaseHealth() DatabaseHealth {
	start := time.Now()

	sqlDB, err := mc.DB.DB()
	if err != nil {
		return DatabaseHealth{
			Status:  "unhealthy",
			Latency: time.Since(start),
		}
	}

	err = sqlDB.Ping()
	latency := time.Since(start)

	if err != nil {
		return DatabaseHealth{
			Status:  "unhealthy",
			Latency: latency,
		}
	}

	stats := sqlDB.Stats()

	return DatabaseHealth{
		Status:      "healthy",
		Connections: stats.OpenConnections,
		Latency:     latency,
	}
}

func (mc *MonitoringController) getSystemHealth() SystemHealth {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return SystemHealth{
		Memory: MemoryStats{
			Allocated: m.Alloc / 1024 / 1024, // Convert to MB
			System:    m.Sys / 1024 / 1024,   // Convert to MB
			GCRuns:    m.NumGC,
		},
		Goroutines: runtime.NumGoroutine(),
		Uptime:     time.Since(startTime).String(),
	}
}

// ReadinessCheck for Kubernetes readiness probe
func (mc *MonitoringController) ReadinessCheck(c *gin.Context) {
	// Check if database is ready
	sqlDB, err := mc.DB.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	err = sqlDB.Ping()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck for Kubernetes liveness probe
func (mc *MonitoringController) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}
