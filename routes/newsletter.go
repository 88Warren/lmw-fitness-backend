package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterNewsletterRoutes(router *gin.Engine, nc *controllers.NewsletterController) {
	newsletterGroup := router.Group("/api/newsletter")

	{
		newsletterGroup.POST("/subscribe", nc.Subscribe)
	}
}
