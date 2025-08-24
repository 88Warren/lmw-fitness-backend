package routes

import (
	"github.com/88warren/lmw-fitness-backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController) {
	api := router.Group("/api")
	{
		api.GET("/", hc.GetHome)
		api.POST("/contact", hc.HandleContactForm)
	}
}
