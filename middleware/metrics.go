package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Metrics struct {
	RequestCount    map[string]int64
	RequestDuration map[string]time.Duration
	ErrorCount      map[string]int64
}

var globalMetrics = &Metrics{
	RequestCount:    make(map[string]int64),
	RequestDuration: make(map[string]time.Duration),
	ErrorCount:      make(map[string]int64),
}

// MetricsCollectionMiddleware collects basic metrics
func MetricsCollectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		path := c.FullPath()
		method := c.Request.Method
		status := c.Writer.Status()

		// Increment request count
		key := method + " " + path
		globalMetrics.RequestCount[key]++

		// Track duration
		globalMetrics.RequestDuration[key] = duration

		// Track errors
		if status >= 400 {
			errorKey := key + " " + strconv.Itoa(status)
			globalMetrics.ErrorCount[errorKey]++
		}
	}
}

// GetMetrics returns current metrics (for health endpoint)
func GetMetrics() *Metrics {
	return globalMetrics
}

// LogMetricsPeriodically logs metrics every 5 minutes
func LogMetricsPeriodically() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			zap.L().Info("Application Metrics",
				zap.Any("request_counts", globalMetrics.RequestCount),
				zap.Any("error_counts", globalMetrics.ErrorCount),
			)
		}
	}()
}
