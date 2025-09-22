package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(router *gin.Engine, healthController *controllers.HealthController) {
	api := router.Group("/api")
	{
		api.GET("/health/live", healthController.LivenessCheck)
		api.GET("/health/ready", healthController.ReadinessCheck)
	}
}
