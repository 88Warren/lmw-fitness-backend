package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController, healthController *controllers.HealthController) {
	router.GET("/api/", hc.GetHome)
	router.POST("/api/contact", hc.HandleContactForm)
	router.GET("/api/health/live", healthController.LivenessCheck)
	router.GET("/api/health/ready", healthController.ReadinessCheck)
	router.POST("/contact", hc.HandleContactForm)
}
