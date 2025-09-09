package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterMonitoringRoutes(router *gin.Engine, mc *controllers.MonitoringController) {
	monitoring := router.Group("/monitoring")
	{
		monitoring.GET("/health", mc.HealthCheck)
		monitoring.GET("/ready", mc.ReadinessCheck)
		monitoring.GET("/live", mc.LivenessCheck)
	}
}
