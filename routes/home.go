package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController, healthController *controllers.HealthController) {
	// API routes
	api := router.Group("/api")
	{
		api.GET("/", hc.GetHome)
		api.POST("/contact", hc.HandleContactForm)
		api.GET("/health/live", healthController.LivenessCheck)
		api.GET("/health/ready", healthController.ReadinessCheck)
		api.POST("/test", hc.TestEndpoint)
	}

	// Images group
	images := router.Group("/images")
	imagesController := controllers.NewImagesController()
	{
		images.GET("/file/:filename", imagesController.GetImage)
	}
}
