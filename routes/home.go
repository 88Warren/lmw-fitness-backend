package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/laurawarren88/LMW_Fitness/controllers"
)

func RegisterHomeRoutes(router *gin.Engine, hc *controllers.HomeController) {
	router.GET("/api/", hc.GetHome)
}
