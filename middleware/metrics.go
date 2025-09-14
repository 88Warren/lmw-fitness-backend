package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Metrics struct {
	mu              sync.RWMutex
	RequestCount    map[string]int64
	RequestDuration map[string]time.Duration
	ErrorCount      map[string]int64
}

var globalMetrics = &Metrics{
	RequestCount:    make(map[string]int64),
	RequestDuration: make(map[string]time.Duration),
	ErrorCount:      make(map[string]int64),
}

func MetricsCollectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		path := c.FullPath()
		method := c.Request.Method
		status := c.Writer.Status()

		globalMetrics.mu.Lock()
		defer globalMetrics.mu.Unlock()

		key := method + " " + path
		globalMetrics.RequestCount[key]++

		globalMetrics.RequestDuration[key] = duration

		if status >= 400 {
			errorKey := key + " " + strconv.Itoa(status)
			globalMetrics.ErrorCount[errorKey]++
		}
	}
}

func GetMetrics() *Metrics {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	metricsCopy := &Metrics{
		RequestCount:    make(map[string]int64),
		RequestDuration: make(map[string]time.Duration),
		ErrorCount:      make(map[string]int64),
	}

	for k, v := range globalMetrics.RequestCount {
		metricsCopy.RequestCount[k] = v
	}
	for k, v := range globalMetrics.RequestDuration {
		metricsCopy.RequestDuration[k] = v
	}
	for k, v := range globalMetrics.ErrorCount {
		metricsCopy.ErrorCount[k] = v
	}

	return metricsCopy
}

func LogMetricsPeriodically() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {

			metrics := GetMetrics()
			zap.L().Info("Application Metrics",
				zap.Any("request_counts", metrics.RequestCount),
				zap.Any("error_counts", metrics.ErrorCount),
			)
		}
	}()
}
