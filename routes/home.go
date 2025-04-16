package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController) {
	router.GET("/api/", hc.GetHome)
	router.POST("/api/contact", hc.HandleContactForm)
}

func HealthCheckRoutes(router *gin.Engine, hc *controllers.HealthController) {
	router.GET("/api//health/live", hc.LivenessCheck)
	router.GET("/api/health/ready", hc.ReadinessCheck)
}
